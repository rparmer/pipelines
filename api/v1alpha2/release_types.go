package v1alpha2

type ReleaseReference struct {
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent.
	// +kubebuilder:validation:Enum=Kustomization;HelmRelease
	// +required
	Kind string `json:"kind"`

	// Name of the referent.
	// +required
	Name string `json:"name"`

	// Namespace of the referent, defaults to the namespace of the Kubernetes resource object that contains the reference.
	// +optional
	// Namespace string `json:"namespace,omitempty"`
}

// func (s *ReleaseReference) String() string {
// 	if s.Namespace != "" {
// 		return fmt.Sprintf("%s/%s/%s", s.Kind, s.Namespace, s.Name)
// 	}
// 	return fmt.Sprintf("%s/%s", s.Kind, s.Name)
// }
