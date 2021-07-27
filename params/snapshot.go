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
	"fmt"
	"strings"

	crdv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/klog/v2"
)

var (
	SnapshotterSecretParams = SecretParamsMap{
		name:               "Snapshotter",
		secretNameKey:      prefixedSnapshotterSecretNameKey,
		secretNamespaceKey: prefixedSnapshotterSecretNamespaceKey,
	}

	SnapshotterListSecretParams = SecretParamsMap{
		name:               "SnapshotterList",
		secretNameKey:      prefixedSnapshotterListSecretNameKey,
		secretNamespaceKey: prefixedSnapshotterListSecretNamespaceKey,
	}
)

// GetSnapshotSecretReference returns a reference to the secret specified in the given nameTemplate
//  and namespaceTemplate, or an error if the templates are not specified correctly.
// No lookup of the referenced secret is performed, and the secret may or may not exist.
//
// supported tokens for name resolution:
// - ${volumesnapshotcontent.name}
// - ${volumesnapshot.namespace}
// - ${volumesnapshot.name}
//
// supported tokens for namespace resolution:
// - ${volumesnapshotcontent.name}
// - ${volumesnapshot.namespace}
//
// an error is returned in the following situations:
// - the nameTemplate or namespaceTemplate contains a token that cannot be resolved
// - the resolved name is not a valid secret name
// - the resolved namespace is not a valid namespace name
func GetSnapshotSecretReference(secretParams SecretParamsMap, snapshotClassParams map[string]string, snapContentName string, snapshot *crdv1.VolumeSnapshot) (*v1.SecretReference, error) {
	nameTemplate, namespaceTemplate, err := VerifyAndGetSecretNameAndNamespaceTemplate(secretParams, snapshotClassParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get name and namespace template from params: %v", err)
	}

	if nameTemplate == "" && namespaceTemplate == "" {
		return nil, nil
	}

	ref := &v1.SecretReference{}

	// Secret namespace template can make use of the VolumeSnapshotContent name, VolumeSnapshot name or namespace.
	// Note that neither of those things are under the control of the VolumeSnapshot user.
	namespaceParams := map[string]string{"volumesnapshotcontent.name": snapContentName}
	// snapshot may be nil when resolving create/delete snapshot secret names because the
	// snapshot may or may not exist at delete time
	if snapshot != nil {
		namespaceParams["volumesnapshot.namespace"] = snapshot.Namespace
	}

	resolvedNamespace, err := resolveTemplate(namespaceTemplate, namespaceParams)
	if err != nil {
		return nil, fmt.Errorf("error resolving value %q: %v", namespaceTemplate, err)
	}
	klog.V(4).Infof("GetSecretReference namespaceTemplate %s, namespaceParams: %+v, resolved %s", namespaceTemplate, namespaceParams, resolvedNamespace)

	if len(validation.IsDNS1123Label(resolvedNamespace)) > 0 {
		if namespaceTemplate != resolvedNamespace {
			return nil, fmt.Errorf("%q resolved to %q which is not a valid namespace name", namespaceTemplate, resolvedNamespace)
		}
		return nil, fmt.Errorf("%q is not a valid namespace name", namespaceTemplate)
	}
	ref.Namespace = resolvedNamespace

	// Secret name template can make use of the VolumeSnapshotContent name, VolumeSnapshot name or namespace.
	// Note that VolumeSnapshot name and namespace are under the VolumeSnapshot user's control.
	nameParams := map[string]string{"volumesnapshotcontent.name": snapContentName}
	if snapshot != nil {
		nameParams["volumesnapshot.name"] = snapshot.Name
		nameParams["volumesnapshot.namespace"] = snapshot.Namespace
	}
	resolvedName, err := resolveTemplate(nameTemplate, nameParams)
	if err != nil {
		return nil, fmt.Errorf("error resolving value %q: %v", nameTemplate, err)
	}
	if len(validation.IsDNS1123Subdomain(resolvedName)) > 0 {
		if nameTemplate != resolvedName {
			return nil, fmt.Errorf("%q resolved to %q which is not a valid secret name", nameTemplate, resolvedName)
		}
		return nil, fmt.Errorf("%q is not a valid secret name", nameTemplate)
	}
	ref.Name = resolvedName

	klog.V(4).Infof("GetSecretReference validated Secret: %+v", ref)
	return ref, nil
}

func RemovePrefixedSnapshotParameters(param map[string]string) (map[string]string, error) {
	newParam := map[string]string{}
	for k, v := range param {
		if strings.HasPrefix(k, csiParameterPrefix) {
			// Check if its well known
			switch k {
			case prefixedSnapshotterSecretNameKey:
			case prefixedSnapshotterSecretNamespaceKey:
			case prefixedSnapshotterListSecretNameKey:
			case prefixedSnapshotterListSecretNamespaceKey:
			default:
				return map[string]string{}, fmt.Errorf("found unknown parameter key \"%s\" with reserved namespace %s", k, csiParameterPrefix)
			}
		} else {
			// Don't strip, add this key-value to new map
			// Deprecated parameters prefixed with "csi" are not stripped to preserve backwards compatibility
			newParam[k] = v
		}
	}
	return newParam, nil
}
