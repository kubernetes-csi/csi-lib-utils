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

const (
	// CSI Parameters prefixed with csiParameterPrefix are not passed through
	// to the driver on CreateVolumeRequest calls. Instead they are intended
	// to used by the CSI external-provisioner and maybe used to populate
	// fields in subsequent CSI calls or Kubernetes API objects.
	csiParameterPrefix = "csi.storage.k8s.io/"

	PrefixedFsTypeKey = csiParameterPrefix + "fstype"

	prefixedDefaultSecretNameKey      = csiParameterPrefix + "secret-name"
	prefixedDefaultSecretNamespaceKey = csiParameterPrefix + "secret-namespace"

	prefixedProvisionerSecretNameKey      = csiParameterPrefix + "provisioner-secret-name"
	prefixedProvisionerSecretNamespaceKey = csiParameterPrefix + "provisioner-secret-namespace"

	prefixedControllerPublishSecretNameKey      = csiParameterPrefix + "controller-publish-secret-name"
	prefixedControllerPublishSecretNamespaceKey = csiParameterPrefix + "controller-publish-secret-namespace"

	prefixedNodeStageSecretNameKey      = csiParameterPrefix + "node-stage-secret-name"
	prefixedNodeStageSecretNamespaceKey = csiParameterPrefix + "node-stage-secret-namespace"

	prefixedNodePublishSecretNameKey      = csiParameterPrefix + "node-publish-secret-name"
	prefixedNodePublishSecretNamespaceKey = csiParameterPrefix + "node-publish-secret-namespace"

	prefixedControllerExpandSecretNameKey      = csiParameterPrefix + "controller-expand-secret-name"
	prefixedControllerExpandSecretNamespaceKey = csiParameterPrefix + "controller-expand-secret-namespace"

	prefixedSnapshotterSecretNameKey      = csiParameterPrefix + "snapshotter-secret-name"
	prefixedSnapshotterSecretNamespaceKey = csiParameterPrefix + "snapshotter-secret-namespace"

	prefixedSnapshotterListSecretNameKey      = csiParameterPrefix + "snapshotter-list-secret-name"
	prefixedSnapshotterListSecretNamespaceKey = csiParameterPrefix + "snapshotter-list-secret-namespace"

	prefixedVolumeSnapshotNameKey        = csiParameterPrefix + "volumesnapshot/name"
	prefixedVolumeSnapshotNamespaceKey   = csiParameterPrefix + "volumesnapshot/namespace"
	prefixedVolumeSnapshotContentNameKey = csiParameterPrefix + "volumesnapshotcontent/name"

	// [Deprecated] CSI Parameters that are put into fields but
	// NOT stripped from the parameters passed to CreateVolume
	provisionerSecretNameKey      = "csiProvisionerSecretName"
	provisionerSecretNamespaceKey = "csiProvisionerSecretNamespace"

	controllerPublishSecretNameKey      = "csiControllerPublishSecretName"
	controllerPublishSecretNamespaceKey = "csiControllerPublishSecretNamespace"

	nodeStageSecretNameKey      = "csiNodeStageSecretName"
	nodeStageSecretNamespaceKey = "csiNodeStageSecretNamespace"

	nodePublishSecretNameKey      = "csiNodePublishSecretName"
	nodePublishSecretNamespaceKey = "csiNodePublishSecretNamespace"

	tokenPVNameKey       = "pv.name"
	tokenPVCNameKey      = "pvc.name"
	tokenPVCNameSpaceKey = "pvc.namespace"
)
