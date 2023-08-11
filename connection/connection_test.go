/*
Copyright 2019 The Kubernetes Authors.

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

package connection

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/container-storage-interface/spec/lib/go/csi"

	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	"k8s.io/component-base/metrics/testutil"
)

func tmpDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "connect")
	require.NoError(t, err, "creating temp directory")
	return dir
}

const (
	serverSock = "server.sock"
)

type identityServer struct{}

func (ids *identityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (ids *identityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (ids *identityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

// startServer creates a gRPC server without any registered services.
// The returned address can be used to connect to it. The cleanup
// function stops it. It can be called multiple times.
func startServer(t *testing.T, tmp string, identity csi.IdentityServer, controller csi.ControllerServer, cmm metrics.CSIMetricsManager) (string, func()) {
	addr := path.Join(tmp, serverSock)
	listener, err := net.Listen("unix", addr)
	require.NoError(t, err, "listening on %s", addr)
	var opts []grpc.ServerOption
	if cmm != nil {
		opts = append(opts,
			grpc.UnaryInterceptor(ExtendedCSIMetricsManager{cmm}.RecordMetricsServerInterceptor),
		)
	}
	server := grpc.NewServer(opts...)
	if identity != nil {
		csi.RegisterIdentityServer(server, identity)
	}
	if controller != nil {
		csi.RegisterControllerServer(server, controller)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(listener); err != nil {
			t.Logf("starting server failed: %s", err)
		}
	}()
	return addr, func() {
		server.Stop()
		wg.Wait()
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			t.Logf("remove Unix socket: %s", err)
		}
	}
}

func TestConnect(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)
	addr, stopServer := startServer(t, tmp, nil, nil, nil)
	defer stopServer()

	conn, err := Connect(addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
	if assert.NoError(t, err, "connect via absolute path") &&
		assert.NotNil(t, conn, "got a connection") {
		assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")
		err = conn.Close()
		assert.NoError(t, err, "closing connection")
	}
}

func TestConnectUnix(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)
	addr, stopServer := startServer(t, tmp, nil, nil, nil)
	defer stopServer()

	conn, err := Connect("unix:///"+addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
	if assert.NoError(t, err, "connect with unix:/// prefix") &&
		assert.NotNil(t, conn, "got a connection") {
		assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")
		err = conn.Close()
		assert.NoError(t, err, "closing connection")
	}
}

func TestConnectWithoutMetrics(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)
	addr, stopServer := startServer(t, tmp, nil, nil, nil)
	defer stopServer()

	conn, err := ConnectWithoutMetrics("unix:///" + addr)
	if assert.NoError(t, err, "connect with unix:/// prefix") &&
		assert.NotNil(t, conn, "got a connection") {
		assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")
		err = conn.Close()
		assert.NoError(t, err, "closing connection")
	}
}

func TestWaitForServer(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)

	// We cannot test that Connect() waits forever for the server
	// to appear, because then we would have to let the test run
	// forever.... What we can test is that it returns shortly
	// after the server appears.
	startTime := time.Now()
	var startTimeServer time.Time
	var stopServer func()
	var wg sync.WaitGroup
	wg.Add(1)
	defer func() {
		wg.Wait()
		stopServer()
	}()
	// Here we pick a relatively long delay before we start the
	// server.  If gRPC did go into an exponential backoff before
	// retrying the connection attempt, then it probably would
	// not react promptly to the server becoming ready. Currently
	// it looks like gRPC tries to connect once per second, with
	// no exponential backoff.
	delay := 10 * time.Second
	go func() {
		defer wg.Done()
		t.Logf("sleeping %s before starting server", delay)
		time.Sleep(delay)
		startTimeServer = time.Now()
		_, stopServer = startServer(t, tmp, nil, nil, nil)
	}()
	conn, err := Connect(path.Join(tmp, serverSock), metrics.NewCSIMetricsManager("fake.csi.driver.io"))
	if assert.NoError(t, err, "connect via absolute path") {
		endTime := time.Now()
		assert.NotNil(t, conn, "got a connection")
		assert.Equal(t, connectivity.Ready.String(), conn.GetState().String(), "connection ready")
		if assert.InEpsilon(t, 1*time.Second, endTime.Sub(startTimeServer), 5, "connection established shortly after server starts") {
			assert.InEpsilon(t, delay, endTime.Sub(startTime), 1)
		}
		err = conn.Close()
		assert.NoError(t, err, "closing connection")
	}
}

func TestTimout(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)

	startTime := time.Now()
	timeout := 5 * time.Second
	conn, err := connect(path.Join(tmp, "no-such.sock"), metrics.NewCSIMetricsManager("fake.csi.driver.io"), []grpc.DialOption{grpc.WithTimeout(timeout)}, nil)
	endTime := time.Now()
	if assert.Error(t, err, "connection should fail") {
		assert.InEpsilon(t, timeout, endTime.Sub(startTime), 1, "connection timeout")
	} else {
		err := conn.Close()
		assert.NoError(t, err, "closing connection")
	}
}

func TestReconnect(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)
	addr, stopServer := startServer(t, tmp, nil, nil, nil)
	defer func() {
		stopServer()
	}()

	// Allow reconnection (the default).
	conn, err := Connect(addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
	if assert.NoError(t, err, "connect via absolute path") &&
		assert.NotNil(t, conn, "got a connection") {
		defer conn.Close()
		assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")

		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unimplemented, errStatus.Code(), "not implemented")
		}

		stopServer()
		startTime := time.Now()
		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			endTime := time.Now()
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unavailable, errStatus.Code(), "connection lost")
			assert.InEpsilon(t, time.Second, endTime.Sub(startTime), 1, "connection loss should be detected quickly")
		}

		// No reconnection either when the server comes back.
		_, stopServer = startServer(t, tmp, nil, nil, nil)
		// We need to give gRPC some time. It does not attempt to reconnect
		// immediately. If we send the method call too soon, the test passes
		// even though a later method call will go through again.
		time.Sleep(5 * time.Second)
		startTime = time.Now()
		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			endTime := time.Now()
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unimplemented, errStatus.Code(), "not implemented")
			assert.InEpsilon(t, time.Second, endTime.Sub(startTime), 1, "connection loss should be covered from quickly")
		}
	}
}

func TestDisconnect(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)
	addr, stopServer := startServer(t, tmp, nil, nil, nil)
	defer func() {
		stopServer()
	}()

	reconnectCount := 0
	conn, err := Connect(addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"), OnConnectionLoss(func() bool {
		reconnectCount++
		// Don't reconnect.
		return false
	}))
	if assert.NoError(t, err, "connect via absolute path") &&
		assert.NotNil(t, conn, "got a connection") {
		defer conn.Close()
		assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")

		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unimplemented, errStatus.Code(), "not implemented")
		}

		stopServer()
		startTime := time.Now()
		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			endTime := time.Now()
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unavailable, errStatus.Code(), "connection lost")
			assert.InEpsilon(t, time.Second, endTime.Sub(startTime), 1, "connection loss should be detected quickly")
		}

		// No reconnection either when the server comes back.
		_, stopServer = startServer(t, tmp, nil, nil, nil)
		// We need to give gRPC some time. It does not attempt to reconnect
		// immediately. If we send the method call too soon, the test passes
		// even though a later method call will go through again.
		time.Sleep(5 * time.Second)
		startTime = time.Now()
		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			endTime := time.Now()
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unavailable, errStatus.Code(), "connection still lost")
			assert.InEpsilon(t, time.Second, endTime.Sub(startTime), 1, "connection loss should be detected quickly")
		}

		assert.Equal(t, 1, reconnectCount, "connection loss callback should be called once")
	}
}

func TestExplicitReconnect(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)
	addr, stopServer := startServer(t, tmp, nil, nil, nil)
	defer func() {
		stopServer()
	}()

	reconnectCount := 0
	conn, err := Connect(addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"), OnConnectionLoss(func() bool {
		reconnectCount++
		// Reconnect.
		return true
	}))
	if assert.NoError(t, err, "connect via absolute path") &&
		assert.NotNil(t, conn, "got a connection") {
		defer conn.Close()
		assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")

		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unimplemented, errStatus.Code(), "not implemented")
		}

		stopServer()
		startTime := time.Now()
		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			endTime := time.Now()
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unavailable, errStatus.Code(), "connection lost")
			assert.InEpsilon(t, time.Second, endTime.Sub(startTime), 1, "connection loss should be detected quickly")
		}

		// No reconnection either when the server comes back.
		_, stopServer = startServer(t, tmp, nil, nil, nil)
		// We need to give gRPC some time. It does not attempt to reconnect
		// immediately. If we send the method call too soon, the test passes
		// even though a later method call will go through again.
		time.Sleep(5 * time.Second)
		startTime = time.Now()
		if err := conn.Invoke(context.Background(), "/connect.v0.Test/Ping", nil, nil); assert.Error(t, err) {
			endTime := time.Now()
			errStatus, _ := status.FromError(err)
			assert.Equal(t, codes.Unimplemented, errStatus.Code(), "connection still lost")
			assert.InEpsilon(t, time.Second, endTime.Sub(startTime), 1, "connection loss should be recovered from quickly")
		}

		assert.Equal(t, 1, reconnectCount, "connection loss callback should be called once")
	}
}

func TestConnectMetrics(t *testing.T) {
	testCases := []struct {
		name            string
		expectedMetrics string
		ctx             context.Context
		checkServer     bool
		cmm             metrics.CSIMetricsManager
	}{
		{
			name: "Regular connection test",
			expectedMetrics: `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
			# TYPE csi_sidecar_operations_seconds histogram
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="0.1"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="0.25"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="0.5"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="1"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="2.5"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="5"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="10"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="15"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="25"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="50"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="120"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="300"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="600"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",le="+Inf"} 1
			csi_sidecar_operations_seconds_sum{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo"} 0
			csi_sidecar_operations_seconds_count{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo"} 1
			`,
			ctx:         context.Background(),
			cmm:         metrics.NewCSIMetricsManager("fake.csi.driver.io"),
			checkServer: true,
		},
		{
			name: "connection test with migrated metrics",
			expectedMetrics: `# HELP csi_sidecar_operations_seconds [ALPHA] Container Storage Interface operation duration with gRPC error code status total
			# TYPE csi_sidecar_operations_seconds histogram
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="0.1"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="0.25"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="0.5"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="1"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="2.5"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="5"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="10"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="15"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="25"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="50"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="120"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="300"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="600"} 1
			csi_sidecar_operations_seconds_bucket{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true",le="+Inf"} 1
			csi_sidecar_operations_seconds_sum{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true"} 0
			csi_sidecar_operations_seconds_count{driver_name="fake.csi.driver.io",grpc_status_code="Unimplemented",method_name="/csi.v1.Identity/GetPluginInfo",migrated="true"} 1
			`,
			ctx: context.WithValue(context.Background(), AdditionalInfoKey, AdditionalInfo{
				Migrated: "true",
			}),
			cmm:         metrics.NewCSIMetricsManagerWithOptions("fake.csi.driver.io", metrics.WithMigration()),
			checkServer: false,
		},
	}
	for _, test := range testCases {
		t.Logf("Running testcase %v", test.name)
		tmp := tmpDir(t)
		defer os.RemoveAll(tmp)
		cmmServer := metrics.NewCSIMetricsManagerForPlugin("fake.csi.driver.io")
		// We have to have a real implementation of the gRPC call, otherwise the metrics
		// interceptor is not called. The CSI identity service is used because it's simple.
		addr, stopServer := startServer(t, tmp, &identityServer{}, nil, cmmServer)
		defer stopServer()

		cmm := test.cmm
		conn, err := Connect(addr, cmm)
		if assert.NoError(t, err, "connect via absolute path") &&
			assert.NotNil(t, conn, "got a connection") {
			defer conn.Close()
			assert.Equal(t, connectivity.Ready, conn.GetState(), "connection ready")

			identityClient := csi.NewIdentityClient(conn)
			if _, err := identityClient.GetPluginInfo(test.ctx, &csi.GetPluginInfoRequest{}); assert.Error(t, err) {
				errStatus, _ := status.FromError(err)
				assert.Equal(t, codes.Unimplemented, errStatus.Code(), "not implemented")
			}
		}

		if err := testutil.GatherAndCompare(
			cmm.GetRegistry(), strings.NewReader(test.expectedMetrics), "csi_sidecar_operations_seconds"); err != nil {
			// Ignore mismatches on csi_sidecar_operations_seconds_sum metric because execution time will vary from test to test.
			err = verifyMetricsError(t, err, "csi_sidecar_operations_seconds_sum")
			if err != nil {
				t.Errorf("Expected client metrics not found -- %v", err)
			}
		}

		if test.checkServer {
			expectedMetrics := strings.Replace(test.expectedMetrics, "csi_sidecar", metrics.SubsystemPlugin, -1)
			if err := testutil.GatherAndCompare(
				cmmServer.GetRegistry(), strings.NewReader(expectedMetrics), "csi_plugin_operations_seconds"); err != nil {
				// Ignore mismatches on csi_sidecar_operations_seconds_sum metric because execution time will vary from test to test.
				err = verifyMetricsError(t, err, metrics.SubsystemPlugin+"_operations_seconds_sum")
				if err != nil {
					t.Errorf("Expected server metrics not found -- %v", err)
				}
			}
		}
	}
}

func verifyMetricsError(t *testing.T, err error, metricToIgnore string) error {
	errStringArr := strings.Split(err.Error(), "got:")

	if len(errStringArr) != 2 {
		return err
	}

	want := errStringArr[0]
	got := strings.TrimSpace(errStringArr[1])

	if want == "" || got == "" {
		return err
	}

	wantArr := strings.Split(err.Error(), "want")
	if len(wantArr) != 2 {
		return err
	}

	want = strings.TrimSpace(wantArr[1])
	want = strings.ReplaceAll(want, "+++ got:", "")

	if matchErr := metrics.VerifyMetricsMatch(want, got, metricToIgnore); matchErr != nil {
		t.Errorf("%v", matchErr)
		return err
	}

	return nil
}
