package main

import (
	"os"

	"github.com/fluxcd/pkg/runtime/client"
	"github.com/fluxcd/pkg/runtime/logger"
	"k8s.io/apimachinery/pkg/runtime"

	helmctrlapi "github.com/fluxcd/helm-controller/api/v2beta1"
	ksctrlapi "github.com/fluxcd/kustomize-controller/api/v1beta2"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/rparmer/pipelines/api/v1alpha1"
	"github.com/rparmer/pipelines/controllers"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	clientgoscheme.AddToScheme(scheme)
	v1alpha1.AddToScheme(scheme)
	ksctrlapi.AddToScheme(scheme)
	helmctrlapi.AddToScheme(scheme)
}

func controller() {
	log := logger.NewLogger(logger.Options{})
	ctrl.SetLogger(log)

	restConfig := client.GetConfigOrDie(client.Options{})
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme: scheme,
		Port:   9443,
		Logger: ctrl.Log,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := (&controllers.PipelineReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Provider")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
