module github.com/weaveworks/profiles

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/fluxcd/helm-controller/api v0.12.0
	github.com/fluxcd/kustomize-controller/api v0.16.0
	github.com/fluxcd/pkg/apis/meta v0.10.1
	github.com/fluxcd/pkg/version v0.1.0
	github.com/fluxcd/source-controller v0.16.0
	github.com/fluxcd/source-controller/api v0.17.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-logr/logr v1.2.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/weaveworks/schemer v0.0.0-20210802122110-338b258ad2ca
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/controller-runtime v0.10.2
)
