package v1alpha2

type ReleaseReference struct {
	// +kubebuilder:validation:Enum=Kustomization;HelmRelease
	Kind string `json:"kind"`
	Name string `json:"name"`
}
