package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineStageSpec struct {
	ReleaseRefs []CrossNamespaceObjectReference `json:"releaseRefs"`
}

// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:object:root=true
type PipelineStage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PipelineStageSpec `json:"spec,omitempty"`
	// Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type PipelineStageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PipelineStage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PipelineStage{}, &PipelineStageList{})
}
