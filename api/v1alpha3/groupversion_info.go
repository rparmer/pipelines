// Package v1alpha1 contains API Schema definitions for the kustomize v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=wego.weave.works
package v1alpha3

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	GroupVersion  = schema.GroupVersion{Group: "wego.weave.works", Version: "v1alpha3"}
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)
