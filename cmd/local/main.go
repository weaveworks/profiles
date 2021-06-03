package main

import (
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/detector"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	scheme := runtime.NewScheme()
	logger := ctrl.Log.WithName("test")
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(sourcev1.AddToScheme(scheme))
	pCatalog := profilesv1.ProfileCatalogSource{
		Spec: profilesv1.ProfileCatalogSourceSpec{
			Repo: "https://github.com/aclevername/profiles-examples",
		},
	}
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	c, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}
	err = detector.Detect(pCatalog, nil, c, logger)
	panic(err)
}
