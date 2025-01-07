/*
Copyright 2018 The Kubernetes Authors.

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

package protosanitizer

import (
	"fmt"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer/test/csitest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// Test case from https://github.com/kubernetes-csi/csi-lib-utils/pull/1#pullrequestreview-180126394.
var testReq = csi.CreateVolumeRequest{
	Name: "test-volume",
	CapacityRange: &csi.CapacityRange{
		RequiredBytes: int64(1024),
		LimitBytes:    int64(1024),
	},
	VolumeCapabilities: []*csi.VolumeCapability{
		{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{
					FsType:     "ext4",
					MountFlags: []string{"flag1", "flag2", "flag3"},
				},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			},
		},
	},
	Secrets:                   map[string]string{"secret1": "secret1", "secret2": "secret2"},
	Parameters:                map[string]string{"param1": "param1", "param2": "param2"},
	VolumeContentSource:       &csi.VolumeContentSource{},
	AccessibilityRequirements: &csi.TopologyRequirement{},
}

func TestStripSecrets(t *testing.T) {
	secretName := "secret-abc"
	secretValue := "123"

	// Current spec.
	createVolume := &csi.CreateVolumeRequest{
		AccessibilityRequirements: &csi.TopologyRequirement{
			Requisite: []*csi.Topology{
				{
					Segments: map[string]string{
						"foo": "bar",
						"x":   "y",
					},
				},
				{
					Segments: map[string]string{
						"a": "b",
					},
				},
			},
		},
		Name: "foo",
		VolumeCapabilities: []*csi.VolumeCapability{
			{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{
						FsType: "ext4",
					},
				},
			},
		},
		CapacityRange: &csi.CapacityRange{
			RequiredBytes: 1024,
		},
		Secrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
	}

	// Revised spec with more secret fields.
	createVolumeFuture := &csitest.CreateVolumeRequest{
		CapacityRange: &csitest.CapacityRange{
			RequiredBytes: 1024,
		},
		MaybeSecretMap: map[int64]*csitest.VolumeCapability{
			1: {ArraySecret: "aaa"},
			2: {ArraySecret: "bbb"},
		},
		Name:         "foo",
		NewSecretInt: 42,
		Seecreets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
		VolumeCapabilities: []*csitest.VolumeCapability{
			{
				AccessType: &csitest.VolumeCapability_Mount{
					Mount: &csitest.VolumeCapability_MountVolume{
						FsType: "ext4",
					},
				},
				ArraySecret: "knock knock",
			},
			{
				ArraySecret: "Who's there?",
			},
		},
		VolumeContentSource: &csitest.VolumeContentSource{
			Type: &csitest.VolumeContentSource_Volume{
				Volume: &csitest.VolumeContentSource_VolumeSource{
					VolumeId:         "abc",
					OneofSecretField: "hello",
				},
			},
			NestedSecretField: "world",
		},
	}

	type testcase struct {
		original, stripped interface{}
	}

	cases := []testcase{
		{nil, "null"},
		{1, "1"},
		{"hello world", `"hello world"`},
		{true, "true"},
		{false, "false"},
		{&csi.CreateVolumeRequest{}, `{}`},
		{&testReq, `{"accessibility_requirements":{},"capacity_range":{"limit_bytes":1024,"required_bytes":1024},"name":"test-volume","parameters":{"param1":"param1","param2":"param2"},"secrets":"***stripped***","volume_capabilities":[{"access_mode":{"mode":"MULTI_NODE_MULTI_WRITER"},"mount":{"fs_type":"ext4","mount_flags":["flag1","flag2","flag3"]}}],"volume_content_source":{}}`},
		{createVolume, `{"accessibility_requirements":{"requisite":[{"segments":{"foo":"bar","x":"y"}},{"segments":{"a":"b"}}]},"capacity_range":{"required_bytes":1024},"name":"foo","secrets":"***stripped***","volume_capabilities":[{"mount":{"fs_type":"ext4"}}]}`},
		{&csitest.CreateVolumeRequest{}, `{}`},
		{createVolumeFuture,
			`{"capacity_range":{"required_bytes":1024},"maybe_secret_map":{"1":{"array_secret":"***stripped***"},"2":{"array_secret":"***stripped***"}},"name":"foo","new_secret_int":"***stripped***","seecreets":"***stripped***","volume_capabilities":[{"array_secret":"***stripped***","mount":{"fs_type":"ext4"}},{"array_secret":"***stripped***"}],"volume_content_source":{"nested_secret_field":"***stripped***","volume":{"oneof_secret_field":"***stripped***","volume_id":"abc"}}}`,
		},
		{&csi.CreateVolumeRequest{
			VolumeCapabilities: []*csi.VolumeCapability{{
				AccessMode: &csi.VolumeCapability_AccessMode{
					// Test for unknown enum value
					Mode: csi.VolumeCapability_AccessMode_Mode(12345),
				},
			}},
		}, `{"volume_capabilities":[{"access_mode":{"mode":12345}}]}`},
	}

	// Message from revised spec as received by a sidecar based on the current spec.
	// The XXX_unrecognized field contains secrets and must not get logged.
	unknownFields := &csi.CreateVolumeRequest{}
	data, err := proto.Marshal(createVolumeFuture)
	if assert.NoError(t, err, "marshall future message") &&
		assert.NoError(t, proto.Unmarshal(data, unknownFields), "unmarshal with unknown fields") {
		cases = append(cases, testcase{unknownFields,
			`{"capacity_range":{"required_bytes":1024},"name":"foo","secrets":"***stripped***","volume_capabilities":[{"mount":{"fs_type":"ext4"}},{}],"volume_content_source":{"volume":{"volume_id":"abc"}}}`,
		})
	}

	for _, c := range cases {
		before := fmt.Sprint(c.original)
		stripped := StripSecrets(c.original)
		if assert.Equal(t, c.stripped, stripped.String(), "unexpected result for fmt s of %s", c.original) {
			if assert.Equal(t, c.stripped, fmt.Sprintf("%v", stripped), "unexpected result for fmt v of %s", c.original) {
				assert.Equal(t, c.stripped, fmt.Sprintf("%+v", stripped), "unexpected result for fmt +v of %s", c.original)
			}
		}
		assert.Equal(t, before, fmt.Sprint(c.original), "original value modified")
	}

	// The secret is hidden because StripSecrets is a struct referencing it.
	dump := fmt.Sprintf("%#v", StripSecrets(createVolume))
	assert.NotContains(t, dump, secretName)
	assert.NotContains(t, dump, secretValue)
}

func BenchmarkStrip(b *testing.B) {
	msg := StripSecrets(&testReq)
	for i := 0; i < b.N; i++ {
		_ = msg.String()
	}
}

func BenchmarkStripLarge(b *testing.B) {
	largeRequest := &csi.CreateVolumeRequest{
		Name: "foo",
		Parameters: map[string]string{
			"param1": "param1",
			"param2": "param2",
		},
		VolumeCapabilities: []*csi.VolumeCapability{
			{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{
						FsType: "ext4",
					},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
			},
		},
		AccessibilityRequirements: &csi.TopologyRequirement{},
	}
	topologies := make([]*csi.Topology, 10000)
	for i := range topologies {
		topologies[i] = &csi.Topology{
			Segments: map[string]string{
				"example.com/instance":          fmt.Sprintf("i-%05d", i),
				"topology.kubernetes.io/zone":   "us-east-1a",
				"topology.kubernetes.io/region": "us-east-1",
			},
		}
	}
	largeRequest.AccessibilityRequirements.Requisite = topologies
	largeRequest.AccessibilityRequirements.Preferred = topologies

	msg := StripSecrets(&largeRequest)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = msg.String()
	}
}
