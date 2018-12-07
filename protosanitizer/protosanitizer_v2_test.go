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
	"strconv"
	"strings"
	"testing"

	csipb "github.com/container-storage-interface/spec/lib/go/csi"
)

func TestSanitizaMsg(t *testing.T) {

	tests := []struct {
		name             string
		msg              *csipb.CreateVolumeRequest
		additionalFields []string
		expected         string
	}{
		{
			name:     "test_with_csi_secret",
			msg:      createVolReq(),
			expected: `name:"test-volume" capacity_range:<required_bytes:1024 limit_bytes:1024 > volume_capabilities:<mount:<fs_type:"ext4" mount_flags:"flag1" mount_flags:"flag2" mount_flags:"flag3" > access_mode:<mode:MULTI_NODE_MULTI_WRITER > > parameters:<key:"param1" value:"param1" > parameters:<key:"param2" value:"param2" > secrets:<key:"secret1" value:"***Sanitized***" > secrets:<key:"secret2" value:"***Sanitized***" > volume_content_source:<volume:<volume_id:"Source_Volume_ID" > > accessibility_requirements:<requisite:<segments:<key:"segment" value:"segment0" > > requisite:<segments:<key:"segment" value:"segment1" > > preferred:<segments:<key:"segment" value:"segment10" > > preferred:<segments:<key:"segment" value:"segment11" > > > `,
		},
		{
			name:             "test_with_csi_secret_and_additional_field",
			msg:              createVolReq(),
			additionalFields: []string{"parameters"},
			expected:         `name:"test-volume" capacity_range:<required_bytes:1024 limit_bytes:1024 > volume_capabilities:<mount:<fs_type:"ext4" mount_flags:"flag1" mount_flags:"flag2" mount_flags:"flag3" > access_mode:<mode:MULTI_NODE_MULTI_WRITER > > parameters:<key:"param1" value:"***Sanitized***" > parameters:<key:"param2" value:"***Sanitized***" > secrets:<key:"secret1" value:"***Sanitized***" > secrets:<key:"secret2" value:"***Sanitized***" > volume_content_source:<volume:<volume_id:"Source_Volume_ID" > > accessibility_requirements:<requisite:<segments:<key:"segment" value:"segment0" > > requisite:<segments:<key:"segment" value:"segment1" > > preferred:<segments:<key:"segment" value:"segment10" > > preferred:<segments:<key:"segment" value:"segment11" > > > `,
		},
	}
	for _, test := range tests {
		result := SanitizeMsg(test.msg, test.additionalFields...)
		if c := strings.Compare(test.expected, result); c != 0 {
			t.Errorf("Test %s failed, expected: \"%s\" got: \"%s\"", test.name, test.expected, result)
		}
	}
}

func rTopo() []*csipb.Topology {
	t := make([]*csipb.Topology, 0)
	for i := 0; i < 2; i++ {
		t = append(t, topo("segment"+strconv.Itoa(i)))
	}
	return t
}

func pTopo() []*csipb.Topology {
	t := make([]*csipb.Topology, 0)
	for i := 10; i < 12; i++ {
		t = append(t, topo("segment"+strconv.Itoa(i)))
	}
	return t
}

func topo(segname string) *csipb.Topology {
	return &csipb.Topology{
		Segments: map[string]string{"segment": segname},
	}
}

func volCap() *csipb.VolumeCapability {
	return &csipb.VolumeCapability{
		AccessType: &csipb.VolumeCapability_Mount{
			Mount: &csipb.VolumeCapability_MountVolume{
				FsType:     "ext4",
				MountFlags: []string{"flag1", "flag2", "flag3"},
			},
		},
		AccessMode: &csipb.VolumeCapability_AccessMode{
			Mode: csipb.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		},
	}
}

func createVolReq() *csipb.CreateVolumeRequest {
	return &csipb.CreateVolumeRequest{
		Name: "test-volume",
		CapacityRange: &csipb.CapacityRange{
			RequiredBytes: int64(1024),
			LimitBytes:    int64(1024),
		},
		VolumeCapabilities: []*csipb.VolumeCapability{volCap()},
		Secrets:            map[string]string{"secret1": "secret1", "secret2": "secret2"},
		Parameters:         map[string]string{"param1": "param1", "param2": "param2"},
		VolumeContentSource: &csipb.VolumeContentSource{
			Type: &csipb.VolumeContentSource_Volume{
				Volume: &csipb.VolumeContentSource_VolumeSource{
					VolumeId: "Source_Volume_ID",
				},
			},
		},
		AccessibilityRequirements: &csipb.TopologyRequirement{
			Requisite: rTopo(),
			Preferred: pTopo(),
		},
	}
}
