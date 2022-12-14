package main

import (
	"context"
	"fmt"
	"strings"

	helmctrlapi "github.com/fluxcd/helm-controller/api/v2beta1"
	ksctrlapi "github.com/fluxcd/kustomize-controller/api/v1beta2"
	pipelineapi "github.com/rparmer/pipelines/api/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func crd() {
	const (
		KustomiztionNameLabel = "kustomize.toolkit.fluxcd.io/name"
	)

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
	// for _, p := range pipelines.Items {
	// 	stages := p.Spec.Stages
	// 	sort.Slice(stages, func(i, j int) bool {
	// 		return stages[i].Order < stages[j].Order
	// 	})
	// }

	// display info for all pipelines
	for i, p := range pipelines.Items {
		// display pipeline info
		fmt.Printf("Pipeline: %s\n", p.Name)
		for _, e := range p.Spec.Environments {
			fmt.Printf("Environment: %s\n", e.Name)
			stage := pipelineapi.PipelineStage{}
			if err := c.Get(context.Background(), client.ObjectKey{Namespace: e.StageRef.Namespace, Name: e.StageRef.Name}, &stage); err != nil {
				panic(err)
			}

			fmt.Printf("Stage: %s\n", stage.Name)
			fmt.Println("Versions:")

			for _, r := range stage.Spec.ReleaseRefs {
				if r.Kind == "Kustomization" {
					k := ksctrlapi.Kustomization{}
					if err := c.Get(context.Background(), client.ObjectKey{Namespace: stage.Namespace, Name: r.Name}, &k); err != nil {
						panic(err)
					}

					ds := appsv1.DeploymentList{}
					if err := c.List(context.Background(), &ds, client.InNamespace(k.Namespace), client.MatchingLabels{
						KustomiztionNameLabel: k.Name,
					}); err != nil {
						panic(err)
					}
					for _, d := range ds.Items {
						fmt.Printf("\tDeployment/%s: ", d.Name)
						for idx, ctr := range d.Spec.Template.Spec.Containers {
							fmt.Printf("%s", strings.Split(ctr.Image, ":")[1])
							if idx < len(d.Spec.Template.Spec.Containers)-1 {
								fmt.Printf(", ")
							}
						}
						fmt.Println()
					}
				} else {
					hr := helmctrlapi.HelmRelease{}
					if err := c.Get(context.Background(), client.ObjectKey{Namespace: stage.Namespace, Name: r.Name}, &hr); err != nil {
						panic(err)
					}
					fmt.Printf("\tHelmRelease/%s: %s\n", hr.Name, hr.Spec.Chart.Spec.Version)
				}
			}
		}
		if i < len(pipelines.Items)-1 {
			fmt.Println()
		}
	}
}
