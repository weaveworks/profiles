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

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/git"
	"github.com/weaveworks/profiles/pkg/profile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ProfileSubscriptionReconciler reconciles a ProfileSubscription object
type ProfileSubscriptionReconciler struct {
	Client client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=weave.works,resources=profilesubscriptions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weave.works,resources=profilesubscriptions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=weave.works,resources=profilesubscriptions/finalizers,verbs=update
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories/status,verbs=get
// +kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases/status,verbs=get

// SetupWithManager sets up the controller with the Manager.
func (r *ProfileSubscriptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ProfileSubscription{}, builder.WithPredicates(
			predicate.GenerationChangedPredicate{},
		)).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ProfileSubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("resource", req.NamespacedName)

	pSub := v1alpha1.ProfileSubscription{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, &pSub)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("resource has been deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get resource")
		return ctrl.Result{}, err
	}

	pDef, err := git.GetProfileDefinition(pSub.Spec.ProfileURL, pSub.Spec.Branch, logger)
	if err != nil {
		r.patchStatusFailing(ctx, &pSub, logger, "error when fetching profile definition")
		return ctrl.Result{}, err
	}

	instance := profile.New(pDef, pSub, r.Client, logger)
	err = instance.CreateArtifacts(ctx)
	if err != nil {
		r.patchStatusFailing(ctx, &pSub, logger, "error when creating profile artifacts")
		return ctrl.Result{}, err
	}

	r.patchStatusRunning(ctx, &pSub, logger)
	// TODO requeuing
	return ctrl.Result{}, nil
}

func (r *ProfileSubscriptionReconciler) patchStatusFailing(ctx context.Context, pSub *v1alpha1.ProfileSubscription, logger logr.Logger, message string) {
	pSub.Status.State = "failing"
	pSub.Status.Message = message
	r.patchStatus(ctx, pSub, logger)
}

func (r *ProfileSubscriptionReconciler) patchStatusRunning(ctx context.Context, pSub *v1alpha1.ProfileSubscription, logger logr.Logger) {
	pSub.Status.State = "running"
	pSub.Status.Message = ""
	r.patchStatus(ctx, pSub, logger)
}

func (r *ProfileSubscriptionReconciler) patchStatus(ctx context.Context, pSub *v1alpha1.ProfileSubscription, logger logr.Logger) {
	key := client.ObjectKeyFromObject(pSub)
	latest := &v1alpha1.ProfileSubscription{}
	if err := r.Client.Get(ctx, key, latest); err != nil {
		logger.Error(err, "failed to get latest resource during patch")
		return
	}
	err := r.Client.Status().Patch(ctx, pSub, client.MergeFrom(latest))
	if err != nil {
		logger.Error(err, "failed to patch status")
	}
}
