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
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// NOTE: Run "make" to regenerate code after modifying this file

// ProfileInstanceSpec defines the desired state of a ProfileInstance
type ProfileInstanceSpec struct {
	// ProfileURL is a fully qualified URL to a profile repo
	ProfileURL string `json:"profileURL,omitempty"`
	// Branch is the git repo branch containing the profile definition (default: main)
	// +kubebuilder:default:=main
	// +optional
	Branch string `json:"branch,omitempty"`

	// Values holds the values for the Helm chart specified in the first artifact
	// +optional
	Values *apiextensionsv1.JSON `json:"values,omitempty"`

	// ValuesFrom holds references to resources containing values for the Helm chart specified in the first artifact
	// +optional
	ValuesFrom []helmv2.ValuesReference `json:"valuesFrom,omitempty"`
}

// ProfileInstanceStatus defines the observed state of ProfileInstance
type ProfileInstanceStatus struct {
	// Conditions holds the conditions for the ProfileInstance
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// ProfileInstance is the Schema for the profileinstances API
type ProfileInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileInstanceSpec   `json:"spec,omitempty"`
	Status ProfileInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileInstanceList contains a list of ProfileInstance
type ProfileInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileInstance{}, &ProfileInstanceList{})
}
