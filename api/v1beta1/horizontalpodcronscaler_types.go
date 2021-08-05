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

package v1beta1

import (
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HorizontalPodCronscalerSpec struct {
	// ScaleTargetRef
	ScaleTargetRef autoscalingv2beta2.CrossVersionObjectReference `json:"scaleTargetRef"`

	// Reference ...
	Reference string `json:"reference,omitempty"`

	// MinReplicas ...
	Replicas int32 `json:"replicas"`

	// Schedule
	Schedule string `json:"schedule"`
}

// ObjectMeta is metadata of objects.
// This is partially copied from metav1.ObjectMeta.
type ObjectMeta struct {
	// Name is the name of the object.
	// +optional
	Name string `json:"name,omitempty"`

	// Labels is a map of string keys and values.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is a map of string keys and values.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// HorizontalPodCronscalerStatus defines the observed state of HorizontalPodCronscaler
type HorizontalPodCronscalerStatus struct {
	// LastTargetTeplicas ...
	LastTargetReplicas *int32 `json:"lastTargetReplicas"`

	// LastSchedule ...
	LastSchedule metav1.Time `json:"lastSchedule"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Reference",type="string",JSONPath=".spec.reference"
//+kubebuilder:printcolumn:name="Schedule",type="string",JSONPath=".spec.schedule"
//+kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
//+kubebuilder:printcolumn:name="Last Schedule",type="date",JSONPath=".status.lastSchedule"

// HorizontalPodCronscaler is the Schema for the horizontalpodcronscalers API
type HorizontalPodCronscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HorizontalPodCronscalerSpec   `json:"spec,omitempty"`
	Status HorizontalPodCronscalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HorizontalPodCronscalerList contains a list of HorizontalPodCronscaler
type HorizontalPodCronscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HorizontalPodCronscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HorizontalPodCronscaler{}, &HorizontalPodCronscalerList{})
}
