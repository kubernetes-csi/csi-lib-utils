/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package secretcache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// fakeClient implements just enough of kubernetes.Interface to see the context passed to Secret reads.
type fakeClient struct {
	kubernetes.Interface
	corev1client.CoreV1Interface

	calls atomic.Int32
	get   func(ctx context.Context, namespace, name string) (*v1.Secret, error)
}

func (f *fakeClient) CoreV1() corev1client.CoreV1Interface { return f }
func (f *fakeClient) Secrets(ns string) corev1client.SecretInterface {
	return &nsClient{f: f, ns: ns}
}

type nsClient struct {
	corev1client.SecretInterface
	f  *fakeClient
	ns string
}

func (n *nsClient) Get(ctx context.Context, name string, _ metav1.GetOptions) (*v1.Secret, error) {
	n.f.calls.Add(1)
	return n.f.get(ctx, n.ns, name)
}

func secret(data map[string]string) *v1.Secret {
	s := &v1.Secret{Data: map[string][]byte{}}
	for k, v := range data {
		s.Data[k] = []byte(v)
	}
	return s
}

func staticClient(data map[string]string) *fakeClient {
	return &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) { return secret(data), nil }}
}

// signal does a non-blocking send so that extra reads fail an assertion instead of hanging the test.
func signal(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *fakeClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

func ref(name string) *v1.SecretReference {
	return &v1.SecretReference{Namespace: "ns", Name: name}
}

func newTestCache(t *testing.T, client *fakeClient, ttl time.Duration, max int) (*Cache, *fakeClock) {
	t.Helper()
	clock := &fakeClock{t: time.Unix(1000, 0)}
	c, err := newCache(context.Background(), client, ttl, max, clock)
	require.NoError(t, err)
	return c, clock
}

func TestHitAndMiss(t *testing.T) {
	client := staticClient(map[string]string{"user": "admin"})
	c, _ := newTestCache(t, client, time.Minute, 10)
	for i := 0; i < 3; i++ {
		got, err := c.Get(context.Background(), ref("a"))
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"user": "admin"}, got)
	}
	assert.EqualValues(t, 1, client.calls.Load())
	_, err := c.Get(context.Background(), ref("b"))
	require.NoError(t, err)
	assert.EqualValues(t, 2, client.calls.Load())
}

func TestNamespaceIsPartOfTheKey(t *testing.T) {
	client := &fakeClient{get: func(_ context.Context, ns, name string) (*v1.Secret, error) {
		return secret(map[string]string{"from": ns + "/" + name}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)
	for i := 0; i < 2; i++ {
		a, err := c.Get(context.Background(), &v1.SecretReference{Namespace: "ns1", Name: "x"})
		require.NoError(t, err)
		b, err := c.Get(context.Background(), &v1.SecretReference{Namespace: "ns2", Name: "x"})
		require.NoError(t, err)
		assert.Equal(t, "ns1/x", a["from"])
		assert.Equal(t, "ns2/x", b["from"])
	}
	assert.EqualValues(t, 2, client.calls.Load())
}

func TestEmptyDataIsNotNil(t *testing.T) {
	client := &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) { return &v1.Secret{}, nil }}
	c, _ := newTestCache(t, client, time.Minute, 10)
	for i := 0; i < 2; i++ {
		got, err := c.Get(context.Background(), ref("a"))
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got)
	}
}

func TestExpiryIsNotExtendedByUse(t *testing.T) {
	client := staticClient(map[string]string{"k": "v"})
	c, clock := newTestCache(t, client, time.Minute, 10)
	ctx := context.Background()
	_, _ = c.Get(ctx, ref("a"))
	clock.Add(59 * time.Second)
	_, _ = c.Get(ctx, ref("a"))
	assert.EqualValues(t, 1, client.calls.Load())
	clock.Add(2 * time.Second)
	_, _ = c.Get(ctx, ref("a"))
	assert.EqualValues(t, 2, client.calls.Load())
}

func TestLRUEviction(t *testing.T) {
	client := staticClient(map[string]string{"k": "v"})
	c, _ := newTestCache(t, client, time.Minute, 2)
	ctx := context.Background()
	for _, n := range []string{"a", "b", "a", "c"} {
		_, err := c.Get(ctx, ref(n))
		require.NoError(t, err)
	}
	assert.EqualValues(t, 3, client.calls.Load())
	_, _ = c.Get(ctx, ref("a"))
	assert.EqualValues(t, 3, client.calls.Load())
	_, _ = c.Get(ctx, ref("b"))
	assert.EqualValues(t, 4, client.calls.Load())
}

func TestErrorsAreNotCached(t *testing.T) {
	var fail atomic.Bool
	fail.Store(true)
	client := &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) {
		if fail.Load() {
			return nil, apierrors.NewNotFound(schema.GroupResource{Resource: "secrets"}, "a")
		}
		return secret(map[string]string{"k": "v"}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)
	_, err := c.Get(context.Background(), ref("a"))
	require.Error(t, err)
	assert.True(t, apierrors.IsNotFound(err), "error must stay inspectable: %v", err)
	assert.Contains(t, err.Error(), "error getting secret a in namespace ns")

	fail.Store(false)
	got, err := c.Get(context.Background(), ref("a"))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"k": "v"}, got)
	assert.EqualValues(t, 2, client.calls.Load())
}

func TestExpiredEntryIsNotUsedInsteadOfError(t *testing.T) {
	var fail atomic.Bool
	client := &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) {
		if fail.Load() {
			return nil, errors.New("api down")
		}
		return secret(map[string]string{"k": "v"}), nil
	}}
	c, clock := newTestCache(t, client, time.Minute, 10)
	_, err := c.Get(context.Background(), ref("a"))
	require.NoError(t, err)
	clock.Add(2 * time.Minute)
	fail.Store(true)
	_, err = c.Get(context.Background(), ref("a"))
	require.Error(t, err)
}

func TestConcurrentCallersShareOneRead(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	client := &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) {
		signal(started)
		<-release
		return secret(map[string]string{"k": "v"}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)

	const n = 50
	var wg sync.WaitGroup
	results := make([]map[string]string, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := c.Get(context.Background(), ref("a"))
			assert.NoError(t, err)
			results[i] = got
		}()
	}
	<-started
	close(release)
	wg.Wait()
	assert.EqualValues(t, 1, client.calls.Load())
	for _, r := range results {
		assert.Equal(t, map[string]string{"k": "v"}, r)
	}
	results[0]["k"] = "changed"
	for _, r := range results[1:] {
		assert.Equal(t, "v", r["k"])
	}
}

func TestDifferentSecretsAreReadConcurrently(t *testing.T) {
	startedA := make(chan struct{}, 1)
	releaseA := make(chan struct{})
	client := &fakeClient{get: func(_ context.Context, _, name string) (*v1.Secret, error) {
		if name == "a" {
			signal(startedA)
			<-releaseA
		}
		return secret(map[string]string{"name": name}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)

	doneA := make(chan error, 1)
	go func() { _, err := c.Get(context.Background(), ref("a")); doneA <- err }()
	<-startedA

	got, err := c.Get(context.Background(), ref("b"))
	require.NoError(t, err)
	assert.Equal(t, "b", got["name"])
	assert.EqualValues(t, 2, client.calls.Load())

	close(releaseA)
	require.NoError(t, <-doneA)
}

func TestFailureIsSharedByWaitersAndNotCached(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	var fail atomic.Bool
	fail.Store(true)
	client := &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) {
		signal(started)
		if fail.Load() {
			<-release
			return nil, errors.New("api down")
		}
		return secret(map[string]string{"k": "v"}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)

	const n = 10
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() { _, err := c.Get(context.Background(), ref("a")); errs <- err }()
	}
	<-started
	close(release)
	for i := 0; i < n; i++ {
		assert.ErrorContains(t, <-errs, "api down")
	}
	callsAfterFailure := client.calls.Load()

	fail.Store(false)
	got, err := c.Get(context.Background(), ref("a"))
	require.NoError(t, err)
	assert.Equal(t, "v", got["k"])
	assert.Greater(t, client.calls.Load(), callsAfterFailure)
}

func TestCancelledCallerDoesNotFailOthers(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	client := &fakeClient{get: func(context.Context, string, string) (*v1.Secret, error) {
		signal(started)
		<-release
		return secret(map[string]string{"k": "v"}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)

	ctxA, cancelA := context.WithCancel(context.Background())
	errA := make(chan error, 1)
	go func() { _, err := c.Get(ctxA, ref("a")); errA <- err }()
	<-started

	resB := make(chan error, 1)
	go func() { _, err := c.Get(context.Background(), ref("a")); resB <- err }()

	cancelA()
	select {
	case err := <-errA:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(10 * time.Second):
		t.Fatal("cancelled caller did not return")
	}

	close(release)
	require.NoError(t, <-resB)
	assert.EqualValues(t, 1, client.calls.Load())
}

func TestAlreadyCancelledContext(t *testing.T) {
	for _, ttl := range []time.Duration{time.Minute, 0, -time.Second} {
		client := staticClient(map[string]string{"k": "v"})
		c, _ := newTestCache(t, client, ttl, 10)
		_, err := c.Get(context.Background(), ref("a"))
		require.NoError(t, err)
		calls := client.calls.Load()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err = c.Get(ctx, ref("a"))
		assert.ErrorIs(t, err, context.Canceled)
		_, err = c.Get(ctx, ref("other"))
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, calls, client.calls.Load(), "ttl %v", ttl)
	}
}

func TestLifetimeCancelStopsSharedRead(t *testing.T) {
	started := make(chan struct{}, 1)
	client := &fakeClient{get: func(ctx context.Context, _, _ string) (*v1.Secret, error) {
		signal(started)
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	lifetime, stop := context.WithCancel(context.Background())
	c, err := New(lifetime, client, time.Minute, 10)
	require.NoError(t, err)

	errc := make(chan error, 1)
	go func() { _, err := c.Get(context.Background(), ref("a")); errc <- err }()
	<-started
	stop()
	select {
	case err := <-errc:
		assert.ErrorIs(t, err, ErrStopped)
		assert.NotErrorIs(t, err, context.Canceled, "must not look like the caller's own cancellation")
	case <-time.After(10 * time.Second):
		t.Fatal("shared read was not cancelled with the lifetime context")
	}

	calls := client.calls.Load()
	_, err = c.Get(context.Background(), ref("b"))
	assert.ErrorIs(t, err, ErrStopped)
	assert.Equal(t, calls, client.calls.Load())
}

func TestLifetimeCancelStillServesCachedEntries(t *testing.T) {
	client := staticClient(map[string]string{"k": "v"})
	lifetime, stop := context.WithCancel(context.Background())
	c, err := New(lifetime, client, time.Minute, 10)
	require.NoError(t, err)
	_, err = c.Get(context.Background(), ref("a"))
	require.NoError(t, err)
	stop()
	got, err := c.Get(context.Background(), ref("a"))
	require.NoError(t, err)
	assert.Equal(t, "v", got["k"])
}

func TestSharedReadHasATimeout(t *testing.T) {
	client := &fakeClient{get: func(ctx context.Context, _, _ string) (*v1.Secret, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)
	c.timeout = 50 * time.Millisecond

	_, err := c.Get(context.Background(), ref("a"))
	require.Error(t, err)
	assert.NotErrorIs(t, err, context.DeadlineExceeded, "must not look like the caller's own deadline")
	assert.NotErrorIs(t, err, ErrStopped)
	assert.Contains(t, err.Error(), "no answer within")
	_, err = c.Get(context.Background(), ref("a"))
	require.Error(t, err)
	assert.EqualValues(t, 2, client.calls.Load())
}

func TestCallerContextDoesNotReachSharedRead(t *testing.T) {
	var readCtx atomic.Value
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	client := &fakeClient{get: func(ctx context.Context, _, _ string) (*v1.Secret, error) {
		readCtx.Store(ctx)
		signal(started)
		<-release
		return secret(nil), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { defer close(done); _, _ = c.Get(ctx, ref("a")) }()
	<-started
	cancel()
	<-done
	assert.NoError(t, readCtx.Load().(context.Context).Err(), "caller cancel must not cancel the shared read")
	close(release)
}

func TestReturnedDataIsACopy(t *testing.T) {
	client := staticClient(map[string]string{"k": "v"})
	c, _ := newTestCache(t, client, time.Minute, 10)
	first, _ := c.Get(context.Background(), ref("a"))
	first["k"] = "changed"
	first["extra"] = "x"
	second, _ := c.Get(context.Background(), ref("a"))
	assert.Equal(t, map[string]string{"k": "v"}, second)
	second["k"] = "changed again"
	delete(second, "k")
	third, _ := c.Get(context.Background(), ref("a"))
	assert.Equal(t, map[string]string{"k": "v"}, third)
}

// The shared read outlives Get, so it must use the key and not the caller's reference.
func TestSharedReadDoesNotLookAtTheCallersReference(t *testing.T) {
	client := &fakeClient{get: func(_ context.Context, _, name string) (*v1.Secret, error) {
		return secret(map[string]string{"name": name}), nil
	}}
	c, _ := newTestCache(t, client, time.Minute, 10)
	r := ref("a")
	c.afterMiss = func() {
		c.afterMiss = nil
		r.Name = "b"
	}
	got, err := c.Get(context.Background(), r)
	require.NoError(t, err)
	assert.Equal(t, "a", got["name"])
}

func TestDisabledReadsEveryTime(t *testing.T) {
	for _, ttl := range []time.Duration{0, -time.Minute} {
		client := staticClient(map[string]string{"k": "v"})
		c, err := New(context.Background(), client, ttl, 0)
		require.NoError(t, err)
		for i := 0; i < 3; i++ {
			_, err := c.Get(context.Background(), ref("a"))
			require.NoError(t, err)
		}
		assert.EqualValues(t, 3, client.calls.Load())
	}
}

func TestNewRejectsBadArguments(t *testing.T) {
	_, err := New(context.Background(), staticClient(nil), time.Minute, 0)
	assert.Error(t, err, "maxEntries")
	_, err = New(nil, staticClient(nil), time.Minute, 1)
	assert.Error(t, err, "nil lifetime")
	_, err = New(nil, staticClient(nil), 0, 1)
	assert.Error(t, err, "nil lifetime, disabled")
	_, err = New(context.Background(), nil, time.Minute, 1)
	assert.Error(t, err, "nil client")
}

func TestNilReference(t *testing.T) {
	client := staticClient(nil)
	for _, ttl := range []time.Duration{0, time.Minute} {
		c, err := New(context.Background(), client, ttl, 1)
		require.NoError(t, err)
		got, err := c.Get(context.Background(), nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	}
	assert.EqualValues(t, 0, client.calls.Load())
}

func TestInvalidReference(t *testing.T) {
	for _, ttl := range []time.Duration{time.Minute, 0} {
		client := staticClient(nil)
		c, _ := newTestCache(t, client, ttl, 10)
		for _, r := range []*v1.SecretReference{
			{Namespace: "", Name: "a"},
			{Namespace: "ns", Name: ""},
			{Namespace: "a/b", Name: "c"},
			{Namespace: "a", Name: "b/c"},
		} {
			_, err := c.Get(context.Background(), r)
			assert.Error(t, err, "%+v", r)
		}
		assert.EqualValues(t, 0, client.calls.Load())
	}
}

// A caller that missed the cache but joins the group after another read finished must not read again.
func TestNoSecondReadAfterConcurrentRead(t *testing.T) {
	client := staticClient(map[string]string{"k": "v"})
	c, _ := newTestCache(t, client, time.Minute, 10)
	c.afterMiss = func() {
		c.afterMiss = nil
		_, err := c.Get(context.Background(), ref("a"))
		assert.NoError(t, err)
	}
	got, err := c.Get(context.Background(), ref("a"))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"k": "v"}, got)
	assert.EqualValues(t, 1, client.calls.Load())
}

func TestRace(t *testing.T) {
	client := staticClient(map[string]string{"k": "v"})
	c, clock := newTestCache(t, client, time.Second, 3)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				name := string(rune('a' + (i+j)%5))
				got, err := c.Get(context.Background(), ref(name))
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, "v", got["k"])
				got["x"] = "y"
				if j%50 == 0 {
					clock.Add(time.Second)
				}
			}
		}()
	}
	wg.Wait()
}
