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

package params

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStripPrefixedCSIParams(t *testing.T) {
	testcases := []struct {
		name           string
		params         map[string]string
		expectedParams map[string]string
		expectErr      bool
	}{
		{
			name:           "no prefix",
			params:         map[string]string{"csiFoo": "bar", "bim": "baz"},
			expectedParams: map[string]string{"csiFoo": "bar", "bim": "baz"},
		},
		{
			name:           "one prefixed",
			params:         map[string]string{prefixedControllerPublishSecretNameKey: "bar", "bim": "baz"},
			expectedParams: map[string]string{"bim": "baz"},
		},
		{
			name:           "prefix in value",
			params:         map[string]string{"foo": prefixedFsTypeKey, "bim": "baz"},
			expectedParams: map[string]string{"foo": prefixedFsTypeKey, "bim": "baz"},
		},
		{
			name: "all known prefixed",
			params: map[string]string{
				prefixedFsTypeKey:                           "csiBar",
				prefixedProvisionerSecretNameKey:            "csiBar",
				prefixedProvisionerSecretNamespaceKey:       "csiBar",
				prefixedControllerPublishSecretNameKey:      "csiBar",
				prefixedControllerPublishSecretNamespaceKey: "csiBar",
				prefixedNodeStageSecretNameKey:              "csiBar",
				prefixedNodeStageSecretNamespaceKey:         "csiBar",
				prefixedNodePublishSecretNameKey:            "csiBar",
				prefixedNodePublishSecretNamespaceKey:       "csiBar",
				prefixedControllerExpandSecretNameKey:       "csiBar",
				prefixedControllerExpandSecretNamespaceKey:  "csiBar",
				prefixedDefaultSecretNameKey:                "csiBar",
				prefixedDefaultSecretNamespaceKey:           "csiBar",
			},
			expectedParams: map[string]string{},
		},
		{
			name: "all known deprecated params not stripped",
			params: map[string]string{
				"fstype":                            "csiBar",
				provisionerSecretNameKey:            "csiBar",
				provisionerSecretNamespaceKey:       "csiBar",
				controllerPublishSecretNameKey:      "csiBar",
				controllerPublishSecretNamespaceKey: "csiBar",
				nodeStageSecretNameKey:              "csiBar",
				nodeStageSecretNamespaceKey:         "csiBar",
				nodePublishSecretNameKey:            "csiBar",
				nodePublishSecretNamespaceKey:       "csiBar",
			},
			expectedParams: map[string]string{
				"fstype":                            "csiBar",
				provisionerSecretNameKey:            "csiBar",
				provisionerSecretNamespaceKey:       "csiBar",
				controllerPublishSecretNameKey:      "csiBar",
				controllerPublishSecretNamespaceKey: "csiBar",
				nodeStageSecretNameKey:              "csiBar",
				nodeStageSecretNamespaceKey:         "csiBar",
				nodePublishSecretNameKey:            "csiBar",
				nodePublishSecretNamespaceKey:       "csiBar",
			},
		},

		{
			name:      "unknown prefixed var",
			params:    map[string]string{csiParameterPrefix + "bim": "baz"},
			expectErr: true,
		},
		{
			name:           "empty",
			params:         map[string]string{},
			expectedParams: map[string]string{},
		},
	}

	for _, tc := range testcases {
		t.Logf("test: %v", tc.name)

		newParams, err := RemovePrefixedParameters(tc.params)
		if err != nil {
			if tc.expectErr {
				continue
			} else {
				t.Fatalf("Encountered unexpected error: %v", err)
			}
		} else {
			if tc.expectErr {
				t.Fatalf("Did not get error when one was expected")
			}
		}

		eq := reflect.DeepEqual(newParams, tc.expectedParams)
		if !eq {
			t.Fatalf("Stripped parameters: %v not equal to expected parameters: %v", newParams, tc.expectedParams)
		}
	}
}

func TestGetVolumeSecretReference(t *testing.T) {
	testcases := map[string]struct {
		secretParams SecretParamsMap
		params       map[string]string
		pvName       string
		pvc          *v1.PersistentVolumeClaim

		expectRef *v1.SecretReference
		expectErr bool
	}{
		"no params": {
			secretParams: NodePublishSecretParams,
			params:       nil,
			expectRef:    nil,
		},
		"empty err": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "", nodePublishSecretNamespaceKey: ""},
			expectErr:    true,
		},
		"[deprecated] name, no namespace": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "foo"},
			expectErr:    true,
		},
		"name, no namespace": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{prefixedNodePublishSecretNameKey: "foo"},
			expectErr:    true,
		},
		"[deprecated] namespace, no name": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNamespaceKey: "foo"},
			expectErr:    true,
		},
		"namespace, no name": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{prefixedNodePublishSecretNamespaceKey: "foo"},
			expectErr:    true,
		},
		"[deprecated] simple - valid": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "name", nodePublishSecretNamespaceKey: "ns"},
			pvc:          &v1.PersistentVolumeClaim{},
			expectRef:    &v1.SecretReference{Name: "name", Namespace: "ns"},
		},
		"deprecated and new both": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "name", nodePublishSecretNamespaceKey: "ns", prefixedNodePublishSecretNameKey: "name", prefixedNodePublishSecretNamespaceKey: "ns"},
			expectErr:    true,
		},
		"deprecated and new names": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "name", nodePublishSecretNamespaceKey: "ns", prefixedNodePublishSecretNameKey: "name"},
			expectErr:    true,
		},
		"deprecated and new namespace": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "name", nodePublishSecretNamespaceKey: "ns", prefixedNodePublishSecretNamespaceKey: "ns"},
			expectErr:    true,
		},
		"deprecated and new mixed": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "name", prefixedNodePublishSecretNamespaceKey: "ns"},
			pvc:          &v1.PersistentVolumeClaim{},
			expectRef:    &v1.SecretReference{Name: "name", Namespace: "ns"},
		},
		"simple - valid": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{prefixedNodePublishSecretNameKey: "name", prefixedNodePublishSecretNamespaceKey: "ns"},
			pvc:          &v1.PersistentVolumeClaim{},
			expectRef:    &v1.SecretReference{Name: "name", Namespace: "ns"},
		},
		"simple - valid, no pvc": {
			secretParams: ProvisionerSecretParams,
			params:       map[string]string{provisionerSecretNameKey: "name", provisionerSecretNamespaceKey: "ns"},
			pvc:          nil,
			expectRef:    &v1.SecretReference{Name: "name", Namespace: "ns"},
		},
		"simple - valid, pvc name and namespace": {
			secretParams: ProvisionerSecretParams,
			params: map[string]string{
				provisionerSecretNameKey:      "param-name",
				provisionerSecretNamespaceKey: "param-ns",
			},
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "ns",
				},
			},
			expectRef: &v1.SecretReference{Name: "param-name", Namespace: "param-ns"},
		},
		"simple - invalid name": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "bad name", nodePublishSecretNamespaceKey: "ns"},
			pvc:          &v1.PersistentVolumeClaim{},
			expectRef:    nil,
			expectErr:    true,
		},
		"simple - invalid namespace": {
			secretParams: NodePublishSecretParams,
			params:       map[string]string{nodePublishSecretNameKey: "name", nodePublishSecretNamespaceKey: "bad ns"},
			pvc:          &v1.PersistentVolumeClaim{},
			expectRef:    nil,
			expectErr:    true,
		},
		"template - PVC name annotations not supported for Provision and Delete": {
			secretParams: ProvisionerSecretParams,
			params: map[string]string{
				prefixedProvisionerSecretNameKey: "static-${pv.name}-${pvc.namespace}-${pvc.name}-${pvc.annotations['akey']}",
			},
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "ns",
				},
			},
			expectErr: true,
		},
		"template - valid nodepublish secret ref": {
			secretParams: NodePublishSecretParams,
			params: map[string]string{
				nodePublishSecretNameKey:      "static-${pv.name}-${pvc.namespace}-${pvc.name}-${pvc.annotations['akey']}",
				nodePublishSecretNamespaceKey: "static-${pv.name}-${pvc.namespace}",
			},
			pvName: "pvname",
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "pvcname",
					Namespace:   "pvcnamespace",
					Annotations: map[string]string{"akey": "avalue"},
				},
			},
			expectRef: &v1.SecretReference{Name: "static-pvname-pvcnamespace-pvcname-avalue", Namespace: "static-pvname-pvcnamespace"},
		},
		"template - valid provisioner secret ref": {
			secretParams: ProvisionerSecretParams,
			params: map[string]string{
				provisionerSecretNameKey:      "static-provisioner-${pv.name}-${pvc.namespace}-${pvc.name}",
				provisionerSecretNamespaceKey: "static-provisioner-${pv.name}-${pvc.namespace}",
			},
			pvName: "pvname",
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pvcname",
					Namespace: "pvcnamespace",
				},
			},
			expectRef: &v1.SecretReference{Name: "static-provisioner-pvname-pvcnamespace-pvcname", Namespace: "static-provisioner-pvname-pvcnamespace"},
		},
		"template - valid, with pvc.name": {
			secretParams: ProvisionerSecretParams,
			params: map[string]string{
				provisionerSecretNameKey:      "${pvc.name}",
				provisionerSecretNamespaceKey: "ns",
			},
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pvcname",
					Namespace: "pvcns",
				},
			},
			expectRef: &v1.SecretReference{Name: "pvcname", Namespace: "ns"},
		},
		"template - valid, provisioner with pvc name and namepsace": {
			secretParams: ProvisionerSecretParams,
			params: map[string]string{
				provisionerSecretNameKey:      "${pvc.name}",
				provisionerSecretNamespaceKey: "${pvc.namespace}",
			},
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pvcname",
					Namespace: "pvcns",
				},
			},
			expectRef: &v1.SecretReference{Name: "pvcname", Namespace: "pvcns"},
		},
		"template - valid, static pvc name and templated namespace": {
			secretParams: ProvisionerSecretParams,
			params: map[string]string{
				provisionerSecretNameKey:      "static-name-1",
				provisionerSecretNamespaceKey: "${pvc.namespace}",
			},
			pvc: &v1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "ns",
				},
			},
			expectRef: &v1.SecretReference{Name: "static-name-1", Namespace: "ns"},
		},
		"template - invalid namespace tokens": {
			secretParams: NodePublishSecretParams,
			params: map[string]string{
				nodePublishSecretNameKey:      "myname",
				nodePublishSecretNamespaceKey: "mynamespace${bar}",
			},
			pvc:       &v1.PersistentVolumeClaim{},
			expectRef: nil,
			expectErr: true,
		},
		"template - invalid name tokens": {
			secretParams: NodePublishSecretParams,
			params: map[string]string{
				nodePublishSecretNameKey:      "myname${foo}",
				nodePublishSecretNamespaceKey: "mynamespace",
			},
			pvc:       &v1.PersistentVolumeClaim{},
			expectRef: nil,
			expectErr: true,
		},
	}

	for k, tc := range testcases {
		t.Run(k, func(t *testing.T) {
			ref, err := GetVolumeSecretReference(tc.secretParams, tc.params, tc.pvName, tc.pvc)
			if err != nil {
				if tc.expectErr {
					return
				}
				t.Fatalf("Did not expect error but got: %v", err)
			} else {
				if tc.expectErr {
					t.Fatalf("Expected error but got none")
				}
			}
			if !reflect.DeepEqual(ref, tc.expectRef) {
				t.Errorf("Expected %v, got %v", tc.expectRef, ref)
			}
		})
	}
}
