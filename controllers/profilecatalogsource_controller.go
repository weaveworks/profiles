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

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
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

func (r *ProfileCatalogSourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("profilecatalogsource", req.NamespacedName)

	pCatalog := v1alpha1.ProfileCatalogSource{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, &pCatalog)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("resource has been deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get resource")
		return ctrl.Result{}, err
	}

	r.Profiles.Add(pCatalog.Name, pCatalog.Spec.Profiles...)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProfileCatalogSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ProfileCatalogSource{}).
		Complete(r)
}
