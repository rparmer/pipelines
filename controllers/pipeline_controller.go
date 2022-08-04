package controllers

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	helmctrlapi "github.com/fluxcd/helm-controller/api/v2beta1"
	ksctrlapi "github.com/fluxcd/kustomize-controller/api/v1beta2"
	v1alpha1 "github.com/rparmer/pipelines/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (p *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Pipeline{}).
		Complete(p)
}

func (p *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, retErr error) {
	pipeline := &v1alpha1.Pipeline{}
	if err := p.Get(ctx, req.NamespacedName, pipeline); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	p.stages(ctx, pipeline)
	return ctrl.Result{}, nil
}

func (p *PipelineReconciler) stages(ctx context.Context, pipeline *v1alpha1.Pipeline) (*v1alpha1.Pipeline, error) {
	stages := make(map[string][]v1alpha1.StageStatus)
	for _, s := range pipeline.Spec.Stages {
		var resources []v1alpha1.StageStatus

		// fetch kustomizations associated to the pipeline
		ks := ksctrlapi.KustomizationList{}
		if err := p.List(context.Background(), &ks, client.InNamespace(s.Namespace), client.MatchingLabels{
			v1alpha1.PipelineNameLabel: pipeline.Name,
		}); err != nil {
			panic(err)
		}

		for _, k := range ks.Items {

			// fetch deployments associated to the kustomization
			ds := appsv1.DeploymentList{}
			if err := p.List(context.Background(), &ds, client.InNamespace(s.Namespace), client.MatchingLabels{
				v1alpha1.KustomiztionNameLabel: k.Name,
			}); err != nil {
				panic(err)
			}
			for _, d := range ds.Items {
				versions := ""
				for idx, ctr := range d.Spec.Template.Spec.Containers {
					versions += strings.Split(ctr.Image, ":")[1]
					if idx < len(d.Spec.Template.Spec.Containers)-1 {
						versions += ", "
					}
				}
				resources = append(resources, v1alpha1.StageStatus{Kind: d.Kind, Name: d.Name, Version: versions})
			}

			// fetch helmreleases associated to the pipeline
			hrs := helmctrlapi.HelmReleaseList{}
			if err := p.List(context.Background(), &hrs, client.InNamespace(s.Namespace), client.MatchingLabels{
				v1alpha1.PipelineNameLabel: pipeline.Name,
			}); err != nil {
				panic(err)
			}
			for _, hr := range hrs.Items {
				resources = append(resources, v1alpha1.StageStatus{Kind: hr.Kind, Name: hr.Name, Version: hr.Spec.Chart.Spec.Version})
			}
		}
		stages[s.Name] = resources
	}

	// store stage info in pipeline status
	pipeline.Status.Stages = stages
	p.patchStatus(ctx, pipeline)

	return pipeline, nil
}

func (p *PipelineReconciler) patchStatus(ctx context.Context, pipeline *v1alpha1.Pipeline) error {
	key := client.ObjectKeyFromObject(pipeline)
	latest := &v1alpha1.Pipeline{}
	if err := p.Client.Get(ctx, key, latest); err != nil {
		return err
	}
	return p.Client.Status().Patch(ctx, pipeline, client.MergeFrom(latest))
}
