package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	helmctrlapi "github.com/fluxcd/helm-controller/api/v2beta1"
	ksctrlapi "github.com/fluxcd/kustomize-controller/api/v1beta1"
	pipelineapi "github.com/rparmer/pipelines/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func cm() {
	const (
		PipelineNameLabel     = "pipelines.wego.weave.works/name"
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

	cm := corev1.ConfigMap{}
	// hardcode cm config for demo
	if err := c.Get(context.Background(), client.ObjectKey{Namespace: "flux-system", Name: "example-pipeline"}, &cm); err != nil {
		panic(err)
	}

	data := cm.Data["pipelines"]

	b, err := ioutil.ReadAll(bytes.NewReader([]byte(data)))
	if err != nil {
		panic(err)
	}

	// convert cm data to PipelineSpec
	pipeline := &pipelineapi.PipelineSpec{}
	err = yaml.Unmarshal(b, pipeline)
	if err != nil {
		panic(err)
	}

	for _, s := range pipeline.Stages {
		fmt.Printf("Stage: %s\n", s.Name)
		fmt.Println("Versions:")

		// fetch kustomizations associated to the pipeline
		ks := ksctrlapi.KustomizationList{}
		if err := c.List(context.Background(), &ks, client.InNamespace(s.Namespace), client.MatchingLabels{
			PipelineNameLabel: "example-pipeline",
		}); err != nil {
			panic(err)
		}

		for _, k := range ks.Items {
			// fetch deployments associated to the kustomization
			ds := appsv1.DeploymentList{}
			if err := c.List(context.Background(), &ds, client.InNamespace(s.Namespace), client.MatchingLabels{
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
		}

		// fetch helmreleases associated to the pipeline
		hrs := helmctrlapi.HelmReleaseList{}
		if err := c.List(context.Background(), &hrs, client.InNamespace(s.Namespace), client.MatchingLabels{
			PipelineNameLabel: "example-pipeline",
		}); err != nil {
			panic(err)
		}
		for _, hr := range hrs.Items {
			fmt.Printf("\tHelmRelease/%s: %s\n", hr.Name, hr.Spec.Chart.Spec.Version)
		}
	}
}
