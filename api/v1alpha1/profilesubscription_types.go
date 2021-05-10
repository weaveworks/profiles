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

// ProfileSubscriptionSpec defines the desired state of a ProfileSubscription
type ProfileSubscriptionSpec struct {
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

	// Version defines the version of the catalog to get the profile from
	Version string `json:"version,omitempty"`

	// Catalog defines the name of the catalog to get the profile from
	Catalog string `json:"catalog,omitempty"`
}

// ProfileSubscriptionStatus defines the observed state of ProfileSubscription
type ProfileSubscriptionStatus struct {
	// Conditions holds the conditions for the ProfileSubscription
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// ProfileSubscription is the Schema for the profilesubscriptions API
type ProfileSubscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileSubscriptionSpec   `json:"spec,omitempty"`
	Status ProfileSubscriptionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileSubscriptionList contains a list of ProfileSubscription
type ProfileSubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileSubscription `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileSubscription{}, &ProfileSubscriptionList{})
}
