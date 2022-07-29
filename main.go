package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	helmctrlapi "github.com/fluxcd/helm-controller/api/v2beta1"
	ksctrlapi "github.com/fluxcd/kustomize-controller/api/v1beta2"
	pipelineapi "github.com/rparmer/pipelines/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	PipelineNameLabel = "pipelines.wego.weave.works/name"
)

func main() {
	cfg := config.GetConfigOrDie()
	c, err := client.New(cfg, client.Options{})
	if err != nil {
		panic(err)
	}

	if err := pipelineapi.AddToScheme(c.Scheme()); err != nil {
		panic(err)
	}
	if err := ksctrlapi.AddToScheme(c.Scheme()); err != nil {
		panic(err)
	}
	if err := helmctrlapi.AddToScheme(c.Scheme()); err != nil {
		panic(err)
	}

	// fetch all pipelines
	pipelines := pipelineapi.PipelineList{}
	if err := c.List(context.Background(), &pipelines); err != nil {
		panic(err)
	}

	// sort pipeline stage order
	for _, p := range pipelines.Items {
		stages := p.Spec.Stages
		sort.Slice(stages, func(i, j int) bool {
			return stages[i].Order < stages[j].Order
		})
	}

	// display info for all pipelines
	for i, p := range pipelines.Items {

		// display pipeline info
		fmt.Printf("Pipeline: %s\n", p.Name)
		for _, s := range p.Spec.Stages {
			fmt.Printf("Stage: %s\n", s.Name)
			fmt.Println("Versions:")

			// fetch deployments associated to the pipeline
			deploys := appsv1.DeploymentList{}
			if err := c.List(context.Background(), &deploys, client.InNamespace(s.Namespace), client.MatchingLabels{
				PipelineNameLabel: p.Name,
			}); err != nil {
				panic(err)
			}
			for _, d := range deploys.Items {
				fmt.Printf("\tDeployment/%s: ", d.Name)
				for idx, ctr := range d.Spec.Template.Spec.Containers {
					fmt.Printf("%s", strings.Split(ctr.Image, ":")[1])
					if idx < len(d.Spec.Template.Spec.Containers)-1 {
						fmt.Printf(", ")
					}
				}
				fmt.Println()
			}

			// fetch helmreleases associated to the pipeline
			hrs := helmctrlapi.HelmReleaseList{}
			if err := c.List(context.Background(), &hrs, client.InNamespace(s.Namespace), client.MatchingLabels{
				PipelineNameLabel: p.Name,
			}); err != nil {
				panic(err)
			}
			for _, hr := range hrs.Items {
				fmt.Printf("\tHelmRelease/%s: %s\n", hr.Name, hr.Spec.Chart.Spec.Version)
			}
		}
		if i < len(pipelines.Items)-1 {
			fmt.Println()
		}
	}
}
