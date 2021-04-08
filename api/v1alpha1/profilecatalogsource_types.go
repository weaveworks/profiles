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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProfileCatalogSourceSpec defines the desired state of ProfileCatalogSource
type ProfileCatalogSourceSpec struct {
	// Profiles is the list of profiles exposed by the catalog
	Profiles []ProfileDescription `json:"profiles,omitempty"`
}

type ProfileDescription struct {
	// Profile name
	Name string `json:"name,omitempty"`
	// Profile description
	Description string `json:"description,omitempty"`
	// Version
	// +optional
	Version string `json:"version,omitempty"`
	// Catalog is the name of the catalog the profile is listed in
	// +optional
	Catalog string `json:"catalog,omitempty"`
	// URL is the full URL path to the profile.yaml
	// +optional
	URL string `json:"url,omitempty"`
	// Maintainer is the name of the author(s)
	// +optional
	Maintainer string `json:"maintainer,omitempty"`
	// Prerequisites are a list of dependencies required by the profile
	// +optional
	Prerequisites []string `json:"prerequisites,omitempty"`
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
