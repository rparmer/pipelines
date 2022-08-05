package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineSpec struct {
	Environments []Environment `json:"environments,omitempty"`
}

type Environment struct {
	Name     string         `json:"name"`
	StageRef StageReference `json:"stageRef"`

	// + optional
	Group string `json:"group,omitempty"`

	// + optional
	Cluster string `json:"cluster,omitempty"`
}

// type PipelineStatus struct {
// 	Stages map[string][]StageStatus `json:"stages,omitempty"`
// }

// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:object:root=true
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PipelineSpec `json:"spec,omitempty"`
	// Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
