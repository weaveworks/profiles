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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// NOTE: Run "make" to regenerate code after modifying this file

// ProfileInstallationSpec defines the desired state of a ProfileInstallation
type ProfileInstallationSpec struct {
	// ConfigMap is the name of the configmap to pull helm values from
	// +optional
	ConfigMap string `json:"configMap,omitempty"`

	// Source defines properties of the source of the profile
	Source *Source `json:"source,omitempty"`

	// Catalog defines properties of the catalog reference
	Catalog *Catalog `json:"catalog,omitempty"`
}

// Source defines the location of the profile
type Source struct {
	// ProfileURL is a fully qualified URL to a profile repo
	URL string `json:"url,omitempty"`

	// Branch is the git repo branch containing the profile definition (default: main)
	// +kubebuilder:default:=main
	// +optional
	Branch string `json:"branch,omitempty"`

	// Path is the location in the git repo containing the profile definition
	// +optional
	Path string `json:"path,omitempty"`

	// Tag is the git tag containing the profile definition
	// +optional
	Tag string `json:"tag,omitempty"`
}

// Catalog defines properties of the catalog this profile is from
type Catalog struct {
	// Version defines the version of the catalog to get the profile from
	Version string `json:"version,omitempty"`

	// Catalog defines the name of the catalog to get the profile from
	Catalog string `json:"catalog,omitempty"`

	// Profile defines the name of the profile
	Profile string `json:"profile,omitempty"`
}

// GetProfileVersion constructs a profile version from the catalog description.
func (p *Catalog) GetProfileVersion() string {
	return fmt.Sprintf("%s/%s", p.Catalog, p.Version)
}

// ProfileInstallationStatus defines the observed state of ProfileInstallation
type ProfileInstallationStatus struct {
	// Conditions holds the conditions for the ProfileInstallation
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// ProfileInstallation is the Schema for the profileinstallations API
type ProfileInstallation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileInstallationSpec   `json:"spec,omitempty"`
	Status ProfileInstallationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileInstallationList contains a list of ProfileInstallation
type ProfileInstallationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileInstallation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileInstallation{}, &ProfileInstallationList{})
}
