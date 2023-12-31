/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WpSpec defines the desired state of Wp
type WpSpec struct {
	Version string `json:"version,omitempty"`
}

// WpStatus defines the observed state of Wp
type WpStatus struct {
	Ready   bool   `json:"ready,omitempty"`
	Version string `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Notified",type=string,JSONPath=`.metadata.resourceVersion`

// Wp is the Schema for the wps API
type Wp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WpSpec   `json:"spec,omitempty"`
	Status WpStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WpList contains a list of Wp
type WpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Wp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Wp{}, &WpList{})
}
