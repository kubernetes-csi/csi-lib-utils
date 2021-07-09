/*
Copyright 2021 The Kubernetes Authors.

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

package accessmodes

import (
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	v1 "k8s.io/api/core/v1"
)

func TestToCSIAccessMode(t *testing.T) {
	tests := []struct {
		name                          string
		pvAccessModes                 []v1.PersistentVolumeAccessMode
		expectedCSIAccessMode         csi.VolumeCapability_AccessMode_Mode
		expectError                   bool
		supportsSingleNodeMultiWriter bool
	}{
		{
			name:                  "RWX",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		},
		{
			name:                  "ROX + RWO",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany, v1.ReadWriteOnce},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_UNKNOWN,
			expectError:           true,
		},
		{
			name:                  "ROX + RWOP",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany, v1.ReadWriteOncePod},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_UNKNOWN,
			expectError:           true,
		},
		{
			name:                  "ROX",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
		},
		{
			name:                  "RWO",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
		{
			name:                  "RWOP",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{v1.ReadWriteOncePod},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
		{
			name:                  "empty",
			pvAccessModes:         []v1.PersistentVolumeAccessMode{},
			expectedCSIAccessMode: csi.VolumeCapability_AccessMode_UNKNOWN,
			expectError:           true,
		},
		{
			name:                          "RWX with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			supportsSingleNodeMultiWriter: true,
		},
		{
			name:                          "ROX + RWO with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany, v1.ReadWriteOnce},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_UNKNOWN,
			expectError:                   true,
			supportsSingleNodeMultiWriter: true,
		},
		{
			name:                          "ROX + RWOP with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany, v1.ReadWriteOncePod},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_UNKNOWN,
			expectError:                   true,
			supportsSingleNodeMultiWriter: true,
		},
		{
			name:                          "ROX with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
			supportsSingleNodeMultiWriter: true,
		},
		{
			name:                          "RWO with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
			supportsSingleNodeMultiWriter: true,
		},
		{
			name:                          "RWOP with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{v1.ReadWriteOncePod},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
			supportsSingleNodeMultiWriter: true,
		},
		{
			name:                          "empty with SINGLE_NODE_MULTI_WRITER capable driver",
			pvAccessModes:                 []v1.PersistentVolumeAccessMode{},
			expectedCSIAccessMode:         csi.VolumeCapability_AccessMode_UNKNOWN,
			expectError:                   true,
			supportsSingleNodeMultiWriter: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			csiAccessMode, err := ToCSIAccessMode(test.pvAccessModes, test.supportsSingleNodeMultiWriter)

			if err == nil && test.expectError {
				t.Errorf("test %s: expected error, got none", test.name)
			}
			if err != nil && !test.expectError {
				t.Errorf("test %s: got error: %s", test.name, err)
			}
			if !test.expectError && csiAccessMode != test.expectedCSIAccessMode {
				t.Errorf("test %s: unexpected access mode: %+v", test.name, csiAccessMode)
			}
		})
	}
}
