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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProfileCatalogSourceSpec defines the desired state of ProfileCatalogSource
type ProfileCatalogSourceSpec struct {
	// Profiles is the list of profiles exposed by the catalog
	Profiles []ProfileCatalogEntry `json:"profiles,omitempty"`
}

// ProfileCatalogEntry defines details about a given profile.
type ProfileCatalogEntry struct {
	// Tag
	// +optional
	// +kubebuilder:validation:Pattern=^([a-zA-Z\-]+\/)?(v)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$

	Tag string `json:"tag,omitempty"`
	// CatalogSource is the name of the catalog the profile is listed in
	// +optional
	CatalogSource string `json:"catalogSource,omitempty"`
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

// ProfileCatalogSource is the Schema for the ProfileCatalogSources API
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
