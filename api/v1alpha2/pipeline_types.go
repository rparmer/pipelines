package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineSpec struct {
	Stages []Stage `json:"stages,omitempty"`
}

type Stage struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// +kubebuilder:validation:Minimum=0
	Order int `json:"order"`

	Release ReleaseReference `json:"releaseRef"`
}

// +kubebuilder:storageversion
// +kubebuilder:object:root=true
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PipelineSpec `json:"spec,omitempty"`
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
