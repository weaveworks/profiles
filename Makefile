SHELL := /bin/bash
# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= 0.0.1

# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "preview,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=preview,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="preview,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= controller-bundle:$(VERSION)

# Image URL to use all building/pushing image targets
IMG ?= profiles-controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.DEFAULT_GOAL := help

all: manager ## Build the manager binary

##@ Build

manager: generate fmt vet ## Build manager binary
	go build -o bin/manager main.go

schema:
	go build -o bin/schema cmd/schema/main.go

fmt: ## Run go fmt against code
	go fmt ./...

vet: ## Run go vet against code
	go vet ./...

lint: ## Run lint against code
	golangci-lint run --exclude-use-default=false --timeout=5m0s

##@ Tests

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: lint generate fmt vet manifests test_deps ## Run unit and integration tests
	source hack/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); ginkgo -r --skipPackage acceptance

acceptance: local-env ## Run acceptance tests
	kubectl -n profiles-system port-forward $(shell kubectl -n profiles-system get pods -l control-plane=controller-manager -o jsonpath={.items[0].metadata.name}) 8000:8000 &
	ginkgo -r tests/acceptance/ || echo "to see logs run: kubectl -n profiles-system logs -f $(shell kubectl -n profiles-system get pods -l control-plane=controller-manager -o jsonpath={.items[0].metadata.name}) manager"

# Running the tests requires the some .toolkit.fluxcd.io CRDs
SOURCE_VER ?= v0.9.0
HELM_VER ?= v0.8.1
KUSTOMIZE_VER ?= v0.10.0
TEST_CRDS:=controllers/testdata/crds
test_deps: ## Fetch test dependencies
	mkdir -p ${TEST_CRDS}
	curl -s --fail https://raw.githubusercontent.com/fluxcd/source-controller/${SOURCE_VER}/config/crd/bases/source.toolkit.fluxcd.io_gitrepositories.yaml \
		-o ${TEST_CRDS}/gitrepositories.yaml
	curl -s --fail https://raw.githubusercontent.com/fluxcd/source-controller/${SOURCE_VER}/config/crd/bases/source.toolkit.fluxcd.io_helmrepositories.yaml \
		-o ${TEST_CRDS}/helmrepositories.yaml
	curl -s --fail https://raw.githubusercontent.com/fluxcd/helm-controller/${HELM_VER}/config/crd/bases/helm.toolkit.fluxcd.io_helmreleases.yaml \
		-o ${TEST_CRDS}/helmreleases.yaml
	curl -s --fail https://raw.githubusercontent.com/fluxcd/kustomize-controller/${KUSTOMIZE_VER}/config/crd/bases/kustomize.toolkit.fluxcd.io_kustomizations.yaml \
		-o ${TEST_CRDS}/kustomize.yaml

##@ Development

local-env: docker-build-local kind-up docker-push-local install undeploy deploy ## Create local kind env and deploy controllers
	flux install --components="source-controller,helm-controller,kustomize-controller"

kind-up: ## Create local kind cluster
	./hack/load-kind.sh

kind-down: ## Tear down local kind cluster
	kind delete cluster --name profiles

install: manifests kustomize ## Install CRDs into a cluster
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from a cluster
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	cd config/manager && $(KUSTOMIZE) edit set image weaveworks/profiles-controller=localhost:5000/${IMG}
	$(KUSTOMIZE) build config/prepare | kubectl apply -f -
	echo "waiting for controller to be ready"
	kubectl -n profiles-system wait --for=condition=available deployment profiles-controller-manager --timeout 5m
	kubectl -n profiles-system wait --for=condition=Ready --all pods --timeout 5m

undeploy: ## UnDeploy controller from the configured Kubernetes cluster in ~/.kube/config
	$(KUSTOMIZE) build config/prepare | kubectl delete --ignore-not-found=true -f -

run: generate fmt vet manifests ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

##@ Generation

manifests: controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: manifests ## Generate code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: bundle
bundle: manifests kustomize ## Generate bundle manifests and metadata, then validate generated files.
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

##@ Docker

docker-build-local: ## Builds the local docker image
	docker build -t localhost:5000/${IMG} .

docker-push-local: ## Builds the local docker image
	docker push localhost:5000/${IMG}

# TODO publish image on release
docker-build: test ## Build the docker image
	docker build -t ${IMG} .

docker-push: ## Push the docker image
	docker push ${IMG}

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

##@ Docs

mdtoc: ## Download mdtoc binary if necessary
	mdtoc -inplace README.md
	GO111MODULE=off go get sigs.k8s.io/mdtoc || true

##@ Utilities

generate-protoc: ## Generate the protobuf files
	buf generate

lint-protoc: ## lint the protocol files
	buf lint

.PHONY: help
help:  ## Display this help. Thanks to https://www.thapaliya.com/en/writings/well-documented-makefiles/
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
mkdir -p $(PROJECT_DIR)/bin ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
