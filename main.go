/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"golang.org/x/sync/errgroup"

	"github.com/weaveworks/profiles/pkg/gateway"
	pgrpc "github.com/weaveworks/profiles/pkg/grpc"
	"github.com/weaveworks/profiles/pkg/interrupt"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/controllers"
	"github.com/weaveworks/profiles/pkg/catalog"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(profilesv1.AddToScheme(scheme))
	utilruntime.Must(sourcev1.AddToScheme(scheme))
	utilruntime.Must(helmv2.AddToScheme(scheme))
	utilruntime.Must(kustomizev1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var enableLeaderElection bool
	var metricsAddr, probeAddr, apiAddr, grpcAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&apiAddr, "profiles-api-bind-address", ":8000", "The address the profiles catalog api binds to.")
	flag.StringVar(&grpcAddr, "profiles-grpc-bind-address", ":50051", "The address the profiles catalog grpc server binds to.")

	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Namespace:              "",
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "b061522d.weave.works",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	profileCatalog := catalog.New()

	if err = (&controllers.ProfileCatalogSourceReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("ProfileCatalogSource"),
		Scheme:   mgr.GetScheme(),
		Profiles: profileCatalog,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ProfileCatalogSource")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// create error group handler
	g, _ := errgroup.WithContext(context.Background())

	grpcServer := pgrpc.NewGRPCServer(setupLog, profileCatalog, grpcAddr)
	setupLog.Info(fmt.Sprintf("starting profiles grpc server at %s", grpcAddr))
	g.Go(func() error {
		return grpcServer.Start()
	})

	setupLog.Info(fmt.Sprintf("starting gateway server at: %s", apiAddr))
	gatewayServer := gateway.NewGatewayServer(setupLog, apiAddr, grpcAddr)
	g.Go(func() error {
		return gatewayServer.Start()
	})

	setupLog.Info("starting manager")
	shutdownContext, managerCancelFunc := context.WithCancel(context.Background())
	g.Go(func() error {
		if err := mgr.Start(shutdownContext); err != nil {
			setupLog.Error(err, "problem running manager")
			return err
		}
		return nil
	})

	setupLog.Info("starting interrupt handler")
	handler := interrupt.NewInterruptHandler(setupLog, grpcServer.Stop, gatewayServer.Stop, func() { managerCancelFunc() })
	g.Go(func() error {
		handler.HandleInterrupts()
		return nil
	})

	if err := g.Wait(); err != nil {
		setupLog.Error(err, "error occurred during server procedures")
		os.Exit(1)
	}
}
