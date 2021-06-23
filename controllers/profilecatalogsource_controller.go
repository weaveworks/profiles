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
	"net/http"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
	"github.com/weaveworks/profiles/pkg/git"
	"github.com/weaveworks/profiles/pkg/gitrepository"
	"github.com/weaveworks/profiles/pkg/scanner"
	corev1 "k8s.io/api/core/v1"
)

// ProfileCatalogSourceReconciler reconciles a ProfileCatalogSource object
type ProfileCatalogSourceReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Profiles *catalog.Catalog
}

// +kubebuilder:rbac:groups=weave.works,resources=profilecatalogsources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weave.works,resources=profilecatalogsources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=weave.works,resources=profilecatalogsources/finalizers,verbs=update

// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories/finalizers,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ProfileCatalogSourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("profilecatalogsource", req.NamespacedName)

	pCatalog := profilesv1.ProfileCatalogSource{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, &pCatalog)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("resource has been deleted")
			r.Profiles.Remove(req.Name)
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get resource")
		return ctrl.Result{}, err
	}

	logger.Info("updating catalog entries", "profiles", pCatalog.Spec.Profiles)
	r.Profiles.Update(pCatalog.Name, pCatalog.Spec.Profiles...)

	timeout := time.Minute * 2
	interval := time.Second * 5
	gitRepoManager := gitrepository.NewManager(ctx, pCatalog.Namespace, r.Client, timeout, interval)
	scanner := scanner.New(gitRepoManager, &git.Client{}, http.DefaultClient, logger)
	for _, repo := range pCatalog.Spec.Repos {
		logger.Info("scan repo for profiles", "repo", repo)
		var secret *corev1.Secret
		if repo.SecretRef != nil {
			secret = &corev1.Secret{}
			objectKey := client.ObjectKey{Name: repo.SecretRef.Name, Namespace: req.Namespace}
			if err := r.Client.Get(ctx, objectKey, secret); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to find secret for repo %v: %w", repo, err)
			}
		}

		profiles, err := scanner.ScanRepository(repo, secret)
		if err != nil {
			return ctrl.Result{}, err
		}
		r.Profiles.Update(pCatalog.Name, profiles...)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProfileCatalogSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&profilesv1.ProfileCatalogSource{}).
		Complete(r)
}
