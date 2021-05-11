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

package controllers

import (
	"context"
	"fmt"
	"strings"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
	"github.com/weaveworks/profiles/pkg/git"
	"github.com/weaveworks/profiles/pkg/profile"
)

const (
	readyFalse   = "False"
	readyTrue    = "True"
	readyUnknown = "Unknown"
)

// ProfileSubscriptionReconciler reconciles a ProfileSubscription object
type ProfileSubscriptionReconciler struct {
	Client   client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Profiles *catalog.Catalog
}

// +kubebuilder:rbac:groups=weave.works,resources=profilesubscriptions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weave.works,resources=profilesubscriptions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=weave.works,resources=profilesubscriptions/finalizers,verbs=update
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories/status,verbs=get
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories/status,verbs=get
// +kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases/status,verbs=get
// +kubebuilder:rbac:groups=kustomize.toolkit.fluxcd.io,resources=kustomizations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kustomize.toolkit.fluxcd.io,resources=kustomizations/status,verbs=get

// SetupWithManager sets up the controller with the Manager.
func (r *ProfileSubscriptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&profilesv1.ProfileSubscription{}, builder.WithPredicates(
			predicate.GenerationChangedPredicate{},
		)). // Owns ensures that changes to resources owned by the pSub cause the pSub to get requeued
		Owns(&sourcev1.GitRepository{}).
		Owns(&helmv2.HelmRelease{}).
		Owns(&sourcev1.HelmRepository{}).
		Owns(&kustomizev1.Kustomization{}).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ProfileSubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("resource", req.NamespacedName)

	pSub := profilesv1.ProfileSubscription{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, &pSub)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("resource has been deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get resource")
		return ctrl.Result{}, err
	}

	//TODO add validation for arguments in PCTL
	if pSub.Spec.ProfileCatalogDescription != nil {
		desc := r.Profiles.GetWithVersion(pSub.Spec.ProfileCatalogDescription.Catalog, pSub.Spec.ProfileCatalogDescription.ProfileName, pSub.Spec.ProfileCatalogDescription.Version)
		if desc == nil {
			logger.Error(err, "profile not found in catalog")
			return ctrl.Result{}, err
		}
		pSub.Spec.ProfileURL = desc.URL
		pSub.Spec.Version = pSub.Spec.ProfileCatalogDescription.GetProfileVersion()
	}

	branchOrTag := pSub.Spec.Branch
	path := profile.GetProfilePathFromSpec(pSub.Spec)
	if pSub.Spec.Version != "" {
		branchOrTag = pSub.Spec.Version
	}
	pDef, err := git.GetProfileDefinition(pSub.Spec.ProfileURL, branchOrTag, path, logger)
	if err != nil {
		if err := r.patchStatus(ctx, &pSub, logger, readyFalse, "FetchProfileFailed", "error when fetching profile definition"); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	instance := profile.New(ctx, pDef, pSub, r.Client, logger)

	if err = instance.ReconcileArtifacts(); err != nil {
		if err := r.patchStatus(ctx, &pSub, logger, readyFalse, "CreateFailed", "error when reconciling profile artifacts"); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	artifactStatus, err := instance.ArtifactStatus()
	if err != nil {
		return ctrl.Result{}, err
	}

	if artifactStatus.ResourcesExist && len(artifactStatus.NotReadyConditions) == 0 {
		return ctrl.Result{}, r.patchStatus(ctx, &pSub, logger, readyTrue, "ArtifactsReady", "all artifact resources ready")
	}

	var messages []string
	status := readyUnknown
	for _, condition := range artifactStatus.NotReadyConditions {
		logger.Info(fmt.Sprintf("%s=%s, message:%s", condition.Type, string(condition.Status), condition.Message))
		messages = append(messages, condition.Message)
		if string(condition.Status) == readyFalse {
			status = readyFalse
		}
	}
	if err := r.patchStatus(ctx, &pSub, logger, status, "ArtifactNotReady", strings.Join(messages, ",")); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ProfileSubscriptionReconciler) patchStatus(ctx context.Context, pSub *profilesv1.ProfileSubscription, logger logr.Logger, readyStatus, reason, message string) error {
	pSub.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionStatus(readyStatus),
			Message:            message,
			Reason:             reason,
			LastTransitionTime: metav1.Now(),
		},
	}

	key := client.ObjectKeyFromObject(pSub)
	latest := &profilesv1.ProfileSubscription{}
	if err := r.Client.Get(ctx, key, latest); err != nil {
		logger.Error(err, "failed to get latest resource during patch")
		return err
	}
	err := r.Client.Status().Patch(ctx, pSub, client.MergeFrom(latest))
	if err != nil {
		logger.Error(err, "failed to patch status")
		return err
	}
	return nil
}
