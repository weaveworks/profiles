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

// HelmChartKind defines properties about the underlying helm chart for an artifact.
const HelmChartKind = "HelmChart"

// KustomizeKind defines a kind containing kustomize yaml files for an artifact.
const KustomizeKind = "Kustomize"

// ProfileKind defines the kind of a profile artifact
const ProfileKind = "Profile"

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// NOTE: Run "make" to regenerate code after modifying this file

// ProfileDefinitionSpec defines the desired state of ProfileDefinition
type ProfileDefinitionSpec struct {
	ProfileDescription `json:",inline"`
	// Artifacts is a list of Profile artifacts
	Artifacts []Artifact `json:"artifacts,omitempty"`
}

// ProfileDescription defines details about a given profile.
type ProfileDescription struct {
	// Profile description
	Description string `json:"description,omitempty"`
	// Maintainer is the name of the author(s)
	// +optional
	Maintainer string `json:"maintainer,omitempty"`
	// Prerequisites are a list of dependencies required by the profile
	// +optional
	Prerequisites []string `json:"prerequisites,omitempty"`
}

// Artifact defines a bundled resource of the components for this profile.
type Artifact struct {
	// Name is the name of the Artifact
	Name string `json:"name,omitempty"`
	// DependsOn is an optional field which defines dependency on other artifacts.
	// +optional
	DependsOn []DependsOn `json:"dependsOn,omitempty"`
	// Chart defines properties to access a remote chart.
	// This is an optional value. It is ignored in case Path is defined.
	// +optional
	Chart *Chart `json:"chart,omitempty"`
	// Profiles defines properties to access a remote profile.
	// +optional
	Profile *Profile `json:"profile,omitempty"`
	// Kustomize defines properties to for a kustomize artifact.
	// +optional
	Kustomize *Kustomize `json:"kustomize,omitempty"`
}

// DependsOn defines an optional artifact name on which this artifact depends on.
type DependsOn struct {
	// Name of the artifact to depend on.
	Name string `json:"name"`
}

// Kustomize defines properties to for a kustomize artifact.
type Kustomize struct {
	// Path is the local path to the Artifact in the Profile repo.
	Path string `json:"path,omitempty"`
}

// Chart defines properties to access remote helm charts.
type Chart struct {
	// URL is the URL of the Helm repository containing a Helm chart and possible values
	// +optional
	URL string `json:"url,omitempty"`
	// Name defines the name of the chart at the remote repository
	// +optional
	Name string `json:"name,omitempty"`
	// Version defines the version of the chart at the remote repository
	// +optional
	Version string `json:"version,omitempty"`
	// Path is the local path to the Artifact in the Profile repo.
	// This is an optional value. If defined, it takes precedence over other Chart fields.
	// +optional
	Path string `json:"path,omitempty"`
	// DefaultValues holds the default values for this Helm release Artifact.
	// These can be overridden by the user, but will otherwise apply.
	// +optional
	DefaultValues string `json:"defaultValues,omitempty"`
}

// Profile defines properties for accessing a profile
type Profile struct {
	// Source defines properties of the source of the profile
	Source *Source `json:"source,omitempty"`
}

// ProfileDefinitionStatus defines the observed state of ProfileDefinition
// This is not used
type ProfileDefinitionStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProfileDefinition is the Schema for the profiles API
type ProfileDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileDefinitionSpec   `json:"spec,omitempty"`
	Status ProfileDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileDefinitionList contains a list of ProfileDefinition
type ProfileDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileDefinition{}, &ProfileDefinitionList{})
}
