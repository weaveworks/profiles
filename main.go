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
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"

	"github.com/weaveworks/profiles/pkg/api"
	"github.com/weaveworks/profiles/pkg/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// setup grpc server details
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on address %s: %v", grpcAddr, err)
	}
	grpcSrv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	reflection.Register(grpcSrv)

	// create the catalog grpc server
	catalogGrpcServer := api.NewCatalogAPI(profileCatalog, ctrl.Log.WithName("api"))
	protos.RegisterProfilesServiceServer(grpcSrv, catalogGrpcServer)

	// handle interrupts
	shutdownContext, managerCancelFunc := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(context.Background())

	// serve grpc apis
	setupLog.Info(fmt.Sprintf("starting profiles grpc server at %s", grpcAddr))
	g.Go(func() error {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			setupLog.Error(err, "unable to start grpc api server")
			return err
		}
		return nil
	})

	// setup grpc-gateway to connect to the grpc server
	mux := gruntime.NewServeMux()
	gopts := []grpc.DialOption{grpc.WithInsecure()}
	if err := protos.RegisterProfilesServiceHandlerFromEndpoint(context.Background(), mux, grpcAddr, gopts); err != nil {
		setupLog.Error(err, "failed to register service handler from endpoint")
		os.Exit(1)
	}

	setupLog.Info(fmt.Sprintf("starting profiles grpc-gateway server at %s", apiAddr))
	server := &http.Server{Addr: apiAddr, Handler: mux}
	g.Go(func() error {
		// ignore server is closing error because the server receives that on graceful shutdown.
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			setupLog.Error(err, "unable to start profiles api server")
			return err
		}
		return nil
	})

	setupLog.Info("starting manager")
	g.Go(func() error {
		if err := mgr.Start(shutdownContext); err != nil {
			setupLog.Error(err, "problem running manager")
			return err
		}
		return nil
	})

	g.Go(func() error {
		interruptChannel := make(chan os.Signal, 2)
		signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case <-interruptChannel:
			done := make(chan struct{})
			// start the timer for the shutdown sequence
			go func() {
				select {
				case <-done:
					return
				case <-time.After(15 * time.Second):
					setupLog.Error(errors.New("timeout"), "graceful shutdown timed out... forcing shutdown")
					os.Exit(1)
				}
			}()
			setupLog.Info("received shutdown signal... gracefully terminating servers...")
			// shutdown the manager
			managerCancelFunc()
			// shutdown grpc-gateway server
			serverTimeoutContext, timeout := context.WithTimeout(context.Background(), 10*time.Second)
			defer timeout()
			if err := server.Shutdown(serverTimeoutContext); err != nil {
				setupLog.Error(err, "Failed to gracefully shutdown server... terminating.")
			}
			// shutdown grpc server
			grpcSrv.GracefulStop()
			setupLog.Info("all done. Goodbye.")
		case <-gctx.Done():
			setupLog.Info("closing signal handler")
			return gctx.Err()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			setupLog.Error(err, "context was cancelled")
		} else {
			setupLog.Error(err, "server error detected")
		}
		os.Exit(1)
	}
}
