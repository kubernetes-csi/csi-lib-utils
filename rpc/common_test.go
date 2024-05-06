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

package rpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"reflect"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	"github.com/stretchr/testify/require"

	"k8s.io/klog/v2/ktesting"
	_ "k8s.io/klog/v2/ktesting/init"
)

func tmpDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "connect")
	require.NoError(t, err, "creating temp directory")
	return dir
}

const (
	serverSock = "server.sock"
)

// startServer creates a gRPC server without any registered services.
// The returned address can be used to connect to it. The cleanup
// function stops it. It can be called multiple times.
func startServer(t *testing.T, tmp string, identity csi.IdentityServer, controller csi.ControllerServer, groupCtrl csi.GroupControllerServer) (string, func()) {
	addr := path.Join(tmp, serverSock)
	listener, err := net.Listen("unix", addr)
	require.NoError(t, err, "listening on %s", addr)
	server := grpc.NewServer()
	if identity != nil {
		csi.RegisterIdentityServer(server, identity)
	}
	if controller != nil {
		csi.RegisterControllerServer(server, controller)
	}
	if groupCtrl != nil {
		csi.RegisterGroupControllerServer(server, groupCtrl)
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

func TestGetDriverName(t *testing.T) {
	tests := []struct {
		name        string
		output      *csi.GetPluginInfoResponse
		injectError bool
		expectError bool
	}{
		{
			name: "success",
			output: &csi.GetPluginInfoResponse{
				Name:          "csi/example",
				VendorVersion: "0.2.0",
				Manifest: map[string]string{
					"hello": "world",
				},
			},
			expectError: false,
		},
		{
			name:        "gRPC error",
			output:      nil,
			injectError: true,
			expectError: true,
		},
		{
			name: "empty name",
			output: &csi.GetPluginInfoResponse{
				Name: "",
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := test.output
			var injectedErr error
			if test.injectError {
				injectedErr = fmt.Errorf("mock error")
			}

			tmp := tmpDir(t)
			defer os.RemoveAll(tmp)
			identity := &fakeIdentityServer{
				pluginInfoResponse: out,
				err:                injectedErr,
			}
			addr, stopServer := startServer(t, tmp, identity, nil, nil)
			defer func() {
				stopServer()
			}()

			_, ctx := ktesting.NewTestContext(t)
			conn, err := connection.Connect(ctx, addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
			if err != nil {
				t.Fatalf("Failed to connect to CSI driver: %s", err)
			}

			name, err := GetDriverName(ctx, conn)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Got error: %v", err)
			}
			if err == nil && name != "csi/example" {
				t.Errorf("Got unexpected name: %q", name)
			}
		})
	}
}

func TestGetPluginCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		output             *csi.GetPluginCapabilitiesResponse
		injectError        bool
		expectCapabilities PluginCapabilitySet
		expectError        bool
	}{
		{
			name: "success",
			output: &csi.GetPluginCapabilitiesResponse{
				Capabilities: []*csi.PluginCapability{
					{
						Type: &csi.PluginCapability_Service_{
							Service: &csi.PluginCapability_Service{
								Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
							},
						},
					},
					{
						Type: &csi.PluginCapability_Service_{
							Service: &csi.PluginCapability_Service{
								Type: csi.PluginCapability_Service_UNKNOWN,
							},
						},
					},
				},
			},
			expectCapabilities: PluginCapabilitySet{
				csi.PluginCapability_Service_CONTROLLER_SERVICE: true,
				csi.PluginCapability_Service_UNKNOWN:            true,
			},
			expectError: false,
		},
		{
			name:        "gRPC error",
			output:      nil,
			injectError: true,
			expectError: true,
		},
		{
			name: "no controller service",
			output: &csi.GetPluginCapabilitiesResponse{
				Capabilities: []*csi.PluginCapability{
					{
						Type: &csi.PluginCapability_Service_{
							Service: &csi.PluginCapability_Service{
								Type: csi.PluginCapability_Service_UNKNOWN,
							},
						},
					},
				},
			},
			expectCapabilities: PluginCapabilitySet{
				csi.PluginCapability_Service_UNKNOWN: true,
			},
			expectError: false,
		},
		{
			name: "empty capability",
			output: &csi.GetPluginCapabilitiesResponse{
				Capabilities: []*csi.PluginCapability{
					{
						Type: nil,
					},
				},
			},
			expectCapabilities: PluginCapabilitySet{},
			expectError:        false,
		},
		{
			name: "no capabilities",
			output: &csi.GetPluginCapabilitiesResponse{
				Capabilities: []*csi.PluginCapability{},
			},
			expectCapabilities: PluginCapabilitySet{},
			expectError:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var injectedErr error
			if test.injectError {
				injectedErr = fmt.Errorf("mock error")
			}

			tmp := tmpDir(t)
			defer os.RemoveAll(tmp)
			identity := &fakeIdentityServer{
				getPluginCapabilitiesResponse: test.output,

				// Make code compatible with gofmt 1.10.2 (used by pull-sig-storage-csi-lib-utils-stable)
				// and 1.11.1 (which will be used by new Prow job) via an extra blank line.
				err: injectedErr,
			}
			addr, stopServer := startServer(t, tmp, identity, nil, nil)
			defer func() {
				stopServer()
			}()

			_, ctx := ktesting.NewTestContext(t)
			conn, err := connection.Connect(ctx, addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
			if err != nil {
				t.Fatalf("Failed to connect to CSI driver: %s", err)
			}

			caps, err := GetPluginCapabilities(ctx, conn)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Got error: %v", err)
			}
			if !reflect.DeepEqual(test.expectCapabilities, caps) {
				t.Errorf("expected capabilities %+v, got %+v", test.expectCapabilities, caps)
			}
		})
	}
}

func TestGetControllerCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		output             *csi.ControllerGetCapabilitiesResponse
		injectError        bool
		expectCapabilities ControllerCapabilitySet
		expectError        bool
	}{
		{
			name: "success",
			output: &csi.ControllerGetCapabilitiesResponse{
				Capabilities: []*csi.ControllerServiceCapability{
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
							},
						},
					},
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
							},
						},
					},
				},
			},
			expectCapabilities: ControllerCapabilitySet{
				csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME:     true,
				csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME: true,
			},
			expectError: false,
		},
		{
			name: "supports read only",
			output: &csi.ControllerGetCapabilitiesResponse{
				Capabilities: []*csi.ControllerServiceCapability{
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_PUBLISH_READONLY,
							},
						},
					},
					{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
							},
						},
					},
				},
			},
			expectCapabilities: ControllerCapabilitySet{
				csi.ControllerServiceCapability_RPC_PUBLISH_READONLY:         true,
				csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME: true,
			},
			expectError: false,
		},
		{
			name:        "gRPC error",
			output:      nil,
			injectError: true,
			expectError: true,
		},
		{
			name: "empty capability",
			output: &csi.ControllerGetCapabilitiesResponse{
				Capabilities: []*csi.ControllerServiceCapability{
					{
						Type: nil,
					},
				},
			},
			expectCapabilities: ControllerCapabilitySet{},
			expectError:        false,
		},
		{
			name: "no capabilities",
			output: &csi.ControllerGetCapabilitiesResponse{
				Capabilities: []*csi.ControllerServiceCapability{},
			},
			expectCapabilities: ControllerCapabilitySet{},
			expectError:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var injectedErr error
			if test.injectError {
				injectedErr = fmt.Errorf("mock error")
			}

			tmp := tmpDir(t)
			defer os.RemoveAll(tmp)
			controller := &fakeControllerServer{
				controllerGetCapabilitiesResponse: test.output,

				// Make code compatible with gofmt 1.10.2 (used by pull-sig-storage-csi-lib-utils-stable)
				// and 1.11.1 (which will be used by new Prow job) via an extra blank line.
				err: injectedErr,
			}
			addr, stopServer := startServer(t, tmp, nil, controller, nil)
			defer func() {
				stopServer()
			}()

			_, ctx := ktesting.NewTestContext(t)
			conn, err := connection.Connect(ctx, addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
			if err != nil {
				t.Fatalf("Failed to connect to CSI driver: %s", err)
			}

			caps, err := GetControllerCapabilities(ctx, conn)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Got error: %v", err)
			}
			if !reflect.DeepEqual(test.expectCapabilities, caps) {
				t.Errorf("expected capabilities %+v, got %+v", test.expectCapabilities, caps)
			}
		})
	}
}

func TestGetGroupControllerCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		output             *csi.GroupControllerGetCapabilitiesResponse
		injectError        bool
		expectCapabilities GroupControllerCapabilitySet
		expectError        bool
	}{
		{
			name: "success",
			output: &csi.GroupControllerGetCapabilitiesResponse{
				Capabilities: []*csi.GroupControllerServiceCapability{
					{
						Type: &csi.GroupControllerServiceCapability_Rpc{
							Rpc: &csi.GroupControllerServiceCapability_RPC{
								Type: csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT,
							},
						},
					},
				},
			},
			expectCapabilities: GroupControllerCapabilitySet{
				csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT: true,
			},
			expectError: false,
		},
		{
			name:        "gRPC error",
			output:      nil,
			injectError: true,
			expectError: true,
		},
		{
			name: "empty capability",
			output: &csi.GroupControllerGetCapabilitiesResponse{
				Capabilities: []*csi.GroupControllerServiceCapability{
					{
						Type: nil,
					},
				},
			},
			expectCapabilities: GroupControllerCapabilitySet{},
			expectError:        false,
		},
		{
			name: "no capabilities",
			output: &csi.GroupControllerGetCapabilitiesResponse{
				Capabilities: []*csi.GroupControllerServiceCapability{},
			},
			expectCapabilities: GroupControllerCapabilitySet{},
			expectError:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var injectedErr error
			if test.injectError {
				injectedErr = fmt.Errorf("mock error")
			}

			tmp := tmpDir(t)
			defer os.RemoveAll(tmp)
			groupCtrl := &fakeGroupControllerServer{
				groupControllerGetCapabilitiesResponse: test.output,

				// Make code compatible with gofmt 1.10.2 (used by pull-sig-storage-csi-lib-utils-stable)
				// and 1.11.1 (which will be used by new Prow job) via an extra blank line.
				err: injectedErr,
			}
			addr, stopServer := startServer(t, tmp, nil, nil, groupCtrl)
			defer func() {
				stopServer()
			}()

			_, ctx := ktesting.NewTestContext(t)
			conn, err := connection.Connect(ctx, addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
			if err != nil {
				t.Fatalf("Failed to connect to CSI driver: %s", err)
			}

			caps, err := GetGroupControllerCapabilities(ctx, conn)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Got error: %v", err)
			}
			if !reflect.DeepEqual(test.expectCapabilities, caps) {
				t.Errorf("expected capabilities %+v, got %+v", test.expectCapabilities, caps)
			}
		})
	}
}

func TestProbeForever(t *testing.T) {
	tests := []struct {
		name        string
		probeCalls  []probeCall
		expectError bool
	}{
		{
			name: "success",
			probeCalls: []probeCall{
				{
					response: &csi.ProbeResponse{
						Ready: &wrappers.BoolValue{Value: true},
					},
				},
			},
			expectError: false,
		},
		{
			name: "success with empty Ready field (true is assumed)",
			probeCalls: []probeCall{
				{
					response: &csi.ProbeResponse{
						Ready: nil,
					},
				},
			},
			expectError: false,
		},
		{
			name: "error",
			probeCalls: []probeCall{
				{
					err: fmt.Errorf("mock error"),
				},
			},
			expectError: true,
		},
		{
			name: "timeout + failure",
			probeCalls: []probeCall{
				{
					err: status.Error(codes.DeadlineExceeded, "timeout"),
				},
				{
					err: fmt.Errorf("mock error"),
				},
			},
			expectError: true,
		},
		{
			name: "timeout + success",
			probeCalls: []probeCall{
				{
					err: status.Error(codes.DeadlineExceeded, "timeout"),
				},
				{
					err: status.Error(codes.DeadlineExceeded, "timeout"),
				},
				{
					response: &csi.ProbeResponse{
						Ready: &wrappers.BoolValue{Value: true},
					},
				},
			},
			expectError: false,
		},
		{
			name: "unready + failure",
			probeCalls: []probeCall{
				{
					response: &csi.ProbeResponse{
						Ready: &wrappers.BoolValue{Value: false},
					},
				},
				{
					err: fmt.Errorf("mock error"),
				},
			},
			expectError: true,
		},
		{
			name: "unready + success",
			probeCalls: []probeCall{
				{
					response: &csi.ProbeResponse{
						Ready: &wrappers.BoolValue{Value: false},
					},
				},
				{
					response: &csi.ProbeResponse{
						Ready: &wrappers.BoolValue{Value: false},
					},
				},
				{
					response: &csi.ProbeResponse{
						Ready: &wrappers.BoolValue{Value: true},
					},
				},
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmp := tmpDir(t)
			defer os.RemoveAll(tmp)
			identity := &fakeIdentityServer{
				probeCalls: test.probeCalls,
			}
			addr, stopServer := startServer(t, tmp, identity, nil, nil)
			defer func() {
				stopServer()
			}()

			_, ctx := ktesting.NewTestContext(t)
			conn, err := connection.Connect(ctx, addr, metrics.NewCSIMetricsManager("fake.csi.driver.io"))
			if err != nil {
				t.Fatalf("Failed to connect to CSI driver: %s", err)
			}

			err = ProbeForever(ctx, conn, time.Second)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Got error: %v", err)
			}
			if len(identity.probeCalls) != identity.probeCallCount {
				t.Errorf("Expected %d probe calls, got %d", len(identity.probeCalls), identity.probeCallCount)
			}
		})
	}
}

type fakeIdentityServer struct {
	pluginInfoResponse            *csi.GetPluginInfoResponse
	getPluginCapabilitiesResponse *csi.GetPluginCapabilitiesResponse
	err                           error

	probeCalls     []probeCall
	probeCallCount int
}

type probeCall struct {
	response *csi.ProbeResponse
	err      error
}

var _ csi.IdentityServer = &fakeIdentityServer{}

func (i *fakeIdentityServer) GetPluginCapabilities(context.Context, *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return i.getPluginCapabilitiesResponse, i.err
}

func (i *fakeIdentityServer) GetPluginInfo(context.Context, *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return i.pluginInfoResponse, i.err
}

func (i *fakeIdentityServer) Probe(context.Context, *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	if i.probeCallCount >= len(i.probeCalls) {
		return nil, fmt.Errorf("Unexpected Probe() call")
	}
	call := i.probeCalls[i.probeCallCount]
	i.probeCallCount++
	return call.response, call.err
}

type fakeControllerServer struct {
	controllerGetCapabilitiesResponse *csi.ControllerGetCapabilitiesResponse
	err                               error
}

func (c *fakeControllerServer) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ControllerModifyVolume(ctx context.Context, request *csi.ControllerModifyVolumeRequest) (*csi.ControllerModifyVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

var _ csi.ControllerServer = &fakeControllerServer{}

func (c *fakeControllerServer) CreateVolume(context.Context, *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) DeleteVolume(context.Context, *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ControllerPublishVolume(context.Context, *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ControllerUnpublishVolume(context.Context, *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ValidateVolumeCapabilities(context.Context, *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ListVolumes(context.Context, *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return c.controllerGetCapabilitiesResponse, c.err
}

func (c *fakeControllerServer) CreateSnapshot(context.Context, *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) DeleteSnapshot(context.Context, *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeControllerServer) ControllerExpandVolume(context.Context, *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

type fakeGroupControllerServer struct {
	groupControllerGetCapabilitiesResponse *csi.GroupControllerGetCapabilitiesResponse
	err                                    error
}

var _ csi.GroupControllerServer = &fakeGroupControllerServer{}

func (c *fakeGroupControllerServer) GroupControllerGetCapabilities(context.Context, *csi.GroupControllerGetCapabilitiesRequest) (*csi.GroupControllerGetCapabilitiesResponse, error) {
	return c.groupControllerGetCapabilitiesResponse, c.err
}

func (c *fakeGroupControllerServer) CreateVolumeGroupSnapshot(context.Context, *csi.CreateVolumeGroupSnapshotRequest) (*csi.CreateVolumeGroupSnapshotResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeGroupControllerServer) DeleteVolumeGroupSnapshot(context.Context, *csi.DeleteVolumeGroupSnapshotRequest) (*csi.DeleteVolumeGroupSnapshotResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c *fakeGroupControllerServer) GetVolumeGroupSnapshot(context.Context, *csi.GetVolumeGroupSnapshotRequest) (*csi.GetVolumeGroupSnapshotResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}
