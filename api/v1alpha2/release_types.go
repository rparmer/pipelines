package v1alpha2

type ReleaseReference struct {
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
