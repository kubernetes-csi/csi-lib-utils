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

package csi_test

import (
	"fmt"
	"testing"

	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	csi "github.com/kubernetes-csi/csi-lib-utils/protosanitizer/test/csi03"
	"github.com/stretchr/testify/assert"
)

func TestStripSecrets(t *testing.T) {
	secretName := "secret-abc"
	secretValue := "123"

	// CSI 0.3.0.
	createVolumeCSI03 := &csi.CreateVolumeRequest{
		AccessibilityRequirements: &csi.TopologyRequirement{
			Requisite: []*csi.Topology{
				&csi.Topology{
					Segments: map[string]string{
						"foo": "bar",
						"x":   "y",
					},
				},
				&csi.Topology{
					Segments: map[string]string{
						"a": "b",
					},
				},
			},
		},
		Name: "foo",
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
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
		ControllerCreateSecrets: map[string]string{
			secretName:   secretValue,
			"secret-xyz": "987",
		},
	}

	type testcase struct {
		original, stripped interface{}
	}

	cases := []testcase{
		{createVolumeCSI03, `{"accessibility_requirements":{"requisite":[{"segments":{"foo":"bar","x":"y"}},{"segments":{"a":"b"}}]},"capacity_range":{"required_bytes":1024},"controller_create_secrets":"***stripped***","name":"foo","volume_capabilities":[{"AccessType":{"Mount":{"fs_type":"ext4"}}}]}`},
	}

	for _, c := range cases {
		before := fmt.Sprint(c.original)
		stripped := protosanitizer.StripSecretsCSI03(c.original)
		if assert.Equal(t, c.stripped, stripped.String(), "unexpected result for fmt s of %s", c.original) {
			if assert.Equal(t, c.stripped, fmt.Sprintf("%v", stripped), "unexpected result for fmt v of %s", c.original) {
				assert.Equal(t, c.stripped, fmt.Sprintf("%+v", stripped), "unexpected result for fmt +v of %s", c.original)
			}
		}
		assert.Equal(t, before, fmt.Sprint(c.original), "original value modified")
	}
}
