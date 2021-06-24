/*
Copyright 2021.

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

package v1alpha1

import (
	"strings"

	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProfileCatalogSourceSpec defines the desired state of ProfileCatalogSource
type ProfileCatalogSourceSpec struct {
	// Profiles is the list of profiles exposed by the catalog
	// +optional
	Profiles []ProfileCatalogEntry `json:"profiles,omitempty"`
	// Repos contains a list of repositories to scan for profiles
	// +optional
	Repos []Repository `json:"repositories,omitempty"`
}

// Repository defines the list of repositories to scan for profiles
type Repository struct {
	// URL is the URL of the repository. When using SSH credentials to access
	// must be in format ssh://git@github.com/stefanprodan/podinfo
	// When using username/password must be in format
	// https://github.com/stefanprodan/podinfo
	URL string `json:"url,omitempty"`
	// The secret name containing the Git credentials.
	// For HTTPS repositories the secret must contain username and password
	// fields.
	// For SSH repositories the secret must contain identity, identity.pub and
	// known_hosts fields.
	// +optional
	SecretRef *meta.LocalObjectReference `json:"secretRef,omitempty"`
}

// ProfileCatalogEntry defines details about a given profile.
type ProfileCatalogEntry struct {
	// Tag
	// +optional
	// +kubebuilder:validation:Pattern=^([a-zA-Z\-]+\/)?(v)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$

	Tag string `json:"tag,omitempty"`
	// CatalogSource is the name of the catalog the profile is listed in
	// +optional
	CatalogSource string `json:"catalog,omitempty"`
	// URL is the full URL path to the profile.yaml
	// +optional
	URL                string `json:"url,omitempty"`
	ProfileDescription `json:",inline"`
}

// ProfileCatalogSourceStatus defines the observed state of ProfileCatalogSource
type ProfileCatalogSourceStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProfileCatalogSource is the Schema for the profilecatalogsources API
type ProfileCatalogSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileCatalogSourceSpec   `json:"spec,omitempty"`
	Status ProfileCatalogSourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileCatalogSourceList contains a list of ProfileCatalogSource
type ProfileCatalogSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileCatalogSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileCatalogSource{}, &ProfileCatalogSourceList{})
}

func GetVersionFromTag(tag string) string {
	splitTag := strings.Split(tag, "/")
	if len(splitTag) == 2 {
		return splitTag[1]
	}
	return tag
}
