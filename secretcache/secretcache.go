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
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilcache "k8s.io/apimachinery/pkg/util/cache"
	"k8s.io/client-go/kubernetes"
)

// fetchTimeout bounds one shared read, including time spent in client-side rate limiting.
const fetchTimeout = 60 * time.Second

// ErrStopped is returned by Get when a read is needed after the lifetime context was cancelled.
var ErrStopped = errors.New("secret cache is stopped")

// Cache caches the data of Secrets. Use New to create one.
type Cache struct {
	client kubernetes.Interface
	// lifetime is the context shared reads run on, so one caller giving up does not fail the others.
	lifetime context.Context
	ttl      time.Duration
	timeout  time.Duration
	// lru is nil when the cache is disabled.
	lru   *utilcache.LRUExpireCache
	group singleflight.Group
	// mu guards pending and keeps it in step with group. It also makes storing and Invalidate atomic.
	mu sync.Mutex
	// pending holds the token of the shared read that group currently has for each Secret.
	// Only that read may be stored.
	pending map[types.NamespacedName]*byte
	// afterMiss is a test hook called between a cache miss and joining the group.
	afterMiss func()
	// beforeRead is a test hook called when a shared read starts.
	beforeRead func()
}

// New returns a Cache that reads Secrets with client. Shared reads are cancelled with lifetime.
// A ttl of zero or less disables caching, and maxEntries is then ignored.
func New(lifetime context.Context, client kubernetes.Interface, ttl time.Duration, maxEntries int) (*Cache, error) {
	return newCache(lifetime, client, ttl, maxEntries, nil)
}

func newCache(lifetime context.Context, client kubernetes.Interface, ttl time.Duration, maxEntries int, clock utilcache.Clock) (*Cache, error) {
	if lifetime == nil {
		return nil, errors.New("lifetime context must not be nil")
	}
	if client == nil {
		return nil, errors.New("client must not be nil")
	}
	c := &Cache{client: client, lifetime: lifetime, ttl: ttl, timeout: fetchTimeout, pending: map[types.NamespacedName]*byte{}}
	if ttl <= 0 {
		return c, nil
	}
	if maxEntries < 1 {
		return nil, fmt.Errorf("maxEntries must be at least 1, got %d", maxEntries)
	}
	if clock == nil {
		c.lru = utilcache.NewLRUExpireCache(maxEntries)
	} else {
		c.lru = utilcache.NewLRUExpireCacheWithClock(maxEntries, clock)
	}
	return c, nil
}

// Get returns a copy of the data of the Secret ref points to, or nil if ref is nil.
// With caching enabled, ctx only bounds how long this caller waits; the read itself runs on the
// lifetime context, and a successful result is cached even if every caller gave up, unless
// Invalidate was called for that Secret meanwhile. With caching disabled, the read uses ctx.
func (c *Cache) Get(ctx context.Context, ref *v1.SecretReference) (map[string]string, error) {
	if ref == nil {
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// "/" is rejected so that the group key below is unambiguous.
	if ref.Namespace == "" || ref.Name == "" || strings.Contains(ref.Namespace, "/") || strings.Contains(ref.Name, "/") {
		return nil, fmt.Errorf("invalid secret reference %q in namespace %q", ref.Name, ref.Namespace)
	}
	if c.lru == nil {
		return c.read(ctx, ref.Namespace, ref.Name)
	}
	// The shared read outlives this call, so it must use the key and not ref.
	key := types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}
	if data, ok := c.lookup(key); ok {
		return maps.Clone(data), nil
	}

	if c.afterMiss != nil {
		c.afterMiss()
	}
	// The token is assigned here and not in the read, which may start only after an Invalidate.
	c.mu.Lock()
	// A read started for data that is already cached would call the API if Invalidate came first.
	if data, ok := c.lookup(key); ok {
		c.mu.Unlock()
		return maps.Clone(data), nil
	}
	token := c.pending[key]
	if token == nil {
		// Each token has its own address, which a zero-size type would not guarantee.
		token = new(byte)
		c.pending[key] = token
	}
	ch := c.group.DoChan(key.String(), func() (any, error) {
		if c.beforeRead != nil {
			c.beforeRead()
		}
		defer c.finish(key, token)
		// Another read may have finished since our lookup; do not read again and restart the TTL.
		if data, ok := c.lookup(key); ok {
			return data, nil
		}
		return c.fetchAndStore(key, token)
	})
	c.mu.Unlock()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}
		return maps.Clone(res.Val.(map[string]string)), nil
	}
}

// Invalidate drops the cached data of the Secret ref points to, for example after a CSI driver
// rejected the credentials. A Get called after Invalidate returns does not get data from before it.
// A read already in progress is not cached, but callers already waiting for it still get its result.
//
// The Secret may not have changed yet: unchanged data is then cached again for the ttl.
// Invalidate drops whatever is cached, even data newer than what the caller used, and the next Get
// costs an API read, so back off when retrying. It does nothing if ref is nil or the cache is disabled.
func (c *Cache) Invalidate(ref *v1.SecretReference) {
	if ref == nil || c.lru == nil {
		return
	}
	key := types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.pending, key)
	c.lru.Remove(key)
	c.group.Forget(key.String())
}

// finish retires the shared read that owns token, unless Invalidate already did.
func (c *Cache) finish(key types.NamespacedName, token *byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pending[key] == token {
		delete(c.pending, key)
		c.group.Forget(key.String())
	}
}

// fetchAndStore is the shared read. Its result is shared by all waiters, the cache keeps its own copy.
func (c *Cache) fetchAndStore(key types.NamespacedName, token *byte) (map[string]string, error) {
	if err := c.lifetime.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStopped, err)
	}
	ctx, cancel := context.WithTimeout(c.lifetime, c.timeout)
	defer cancel()
	data, err := c.read(ctx, key.Namespace, key.Name)
	if err != nil {
		// Keep these apart from the cancellation or deadline of a caller.
		if c.lifetime.Err() != nil {
			return nil, fmt.Errorf("%w: %v", ErrStopped, err)
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("error getting secret %s in namespace %s: no answer within %v", key.Name, key.Namespace, c.timeout)
		}
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// Invalidate removes the token; this read may then be older than the change it was called for.
	if c.pending[key] == token {
		c.lru.Add(key, maps.Clone(data), c.ttl)
	}
	return data, nil
}

func (c *Cache) lookup(key types.NamespacedName) (map[string]string, bool) {
	v, ok := c.lru.Get(key)
	if !ok {
		return nil, false
	}
	return v.(map[string]string), true
}

func (c *Cache) read(ctx context.Context, namespace, name string) (map[string]string, error) {
	secret, err := c.client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting secret %s in namespace %s: %w", name, namespace, err)
	}
	data := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		data[k] = string(v)
	}
	return data, nil
}
