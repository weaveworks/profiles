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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// NOTE: Run "make" to regenerate code after modifying this file

// ProfileDefinitionSpec defines the desired state of ProfileDefinition
type ProfileDefinitionSpec struct {
	// Description is some text to allow a user to identify what this profile installs.
	Description string `json:"description,omitempty"`
	// Artifacts is a list of Profile artifacts
	// can be one of HelmChart, TODO
	Artifacts []Artifact `json:"artifacts,omitempty"`
}

type Artifact struct {
	// Name is the name of the Artifact
	Name string `json:"name,omitempty"`
	// Path is the local path to the Artifact in the Profile repo
	Path string `json:"path,omitempty"`
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
