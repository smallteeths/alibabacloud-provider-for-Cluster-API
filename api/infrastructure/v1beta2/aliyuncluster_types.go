/*
Copyright (c) 2024-2025, Alibaba Cloud and its affiliates;

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AliyunClusterSpec defines the desired state of AliyunCluster
type AliyunClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// foo is an example field of AliyunCluster. Edit aliyuncluster_types.go to remove/update
	// +optional
	Foo *string `json:"foo,omitempty"`
}

// AliyunClusterStatus defines the observed state of AliyunCluster.
type AliyunClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AliyunCluster is the Schema for the aliyunclusters API
type AliyunCluster struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of AliyunCluster
	// +required
	Spec AliyunClusterSpec `json:"spec"`

	// status defines the observed state of AliyunCluster
	// +optional
	Status AliyunClusterStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// AliyunClusterList contains a list of AliyunCluster
type AliyunClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AliyunCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AliyunCluster{}, &AliyunClusterList{})
}
