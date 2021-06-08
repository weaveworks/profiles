module github.com/weaveworks/profiles

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/fluxcd/helm-controller/api v0.10.1
	github.com/fluxcd/kustomize-controller/api v0.12.2
	github.com/fluxcd/pkg/version v0.0.1
	github.com/fluxcd/source-controller/api v0.13.1
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.13.0
	golang.org/x/tools v0.1.0 // indirect
	k8s.io/api v0.20.5
	k8s.io/apiextensions-apiserver v0.20.5
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.20.5
	knative.dev/pkg v0.0.0-20210412173742-b51994e3b312
	sigs.k8s.io/controller-runtime v0.8.3
)
