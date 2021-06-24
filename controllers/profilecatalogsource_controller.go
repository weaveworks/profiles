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
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
	"github.com/weaveworks/profiles/pkg/git"
	"github.com/weaveworks/profiles/pkg/gitrepository"
	"github.com/weaveworks/profiles/pkg/scanner"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProfileCatalogSourceReconciler reconciles a ProfileCatalogSource object
type ProfileCatalogSourceReconciler struct {
	client.Client
	log        logr.Logger
	s          *runtime.Scheme
	Profiles   *catalog.Catalog
	newScanner NewScanner
	timeout    time.Duration
	interval   time.Duration
}

func NewCatalogSourceReconciler(c client.Client, log logr.Logger, scheme *runtime.Scheme, profiles *catalog.Catalog) *ProfileCatalogSourceReconciler {
	return &ProfileCatalogSourceReconciler{
		Client:     c,
		log:        log,
		s:          scheme,
		Profiles:   profiles,
		newScanner: scanner.New,
		timeout:    time.Minute * 2,
		interval:   time.Second * 5,
	}
}

type NewScanner func(gitRepositoryManager scanner.GitRepositoryManager, gitClient scanner.GitClient, httpClients scanner.HTTPClient, logger logr.Logger) scanner.RepoScanner

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
	logger := r.log.WithValues("profilecatalogsource", req.NamespacedName)

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

	//can configre spec.Profiles or spec.Repositories, not both.
	if len(pCatalog.Spec.Profiles) > 0 {
		logger.Info("updating catalog entries", "profiles", pCatalog.Spec.Profiles)
		r.Profiles.AddOrReplace(pCatalog.Name, pCatalog.Spec.Profiles...)
		return ctrl.Result{}, nil
	}

	gitRepoManager := gitrepository.NewManager(ctx, pCatalog.Namespace, r.Client, r.timeout, r.interval)
	scanner := r.newScanner(gitRepoManager, &git.Client{}, http.DefaultClient, logger)
	if r.Profiles.CatalogExists(pCatalog.Name) {
		logger.Info("updating catalog", "catalog", pCatalog.Name)
		return ctrl.Result{}, r.updateCatalogWithRepositories(ctx, &pCatalog, scanner, logger)
	}

	logger.Info("creating catalog", "catalog", pCatalog.Name)
	return ctrl.Result{}, r.createCatalogWithRepositories(ctx, &pCatalog, scanner, logger)
}

func (r *ProfileCatalogSourceReconciler) createCatalogWithRepositories(ctx context.Context, pCatalog *profilesv1.ProfileCatalogSource, scanner scanner.RepoScanner, logger logr.Logger) error {
	for _, repo := range pCatalog.Spec.Repos {
		logger.Info("scan repo for profiles", "repo", repo)
		var secret *corev1.Secret
		if repo.SecretRef != nil {
			secret = &corev1.Secret{}
			objectKey := client.ObjectKey{Name: repo.SecretRef.Name, Namespace: pCatalog.Namespace}
			if err := r.Client.Get(ctx, objectKey, secret); err != nil {
				return fmt.Errorf("failed to find secret for repo %v: %w", repo, err)
			}
		}

		profiles, newTags, err := scanner.ScanRepository(repo, secret, nil)
		if err != nil {
			return err
		}

		createOrReplaceScannedRepositoryStatus(pCatalog, repo, newTags)
		logger.Info("updating catalog with scanning reuslts", "profiles", profiles)
		r.Profiles.Append(pCatalog.Name, profiles...)
	}

	if err := r.Client.Status().Update(ctx, pCatalog); err != nil {
		return fmt.Errorf("failed to patch status: %w", err)
	}
	return nil
}

func createOrReplaceScannedRepositoryStatus(pCatalog *profilesv1.ProfileCatalogSource, repo profilesv1.Repository, newTags []string) {
	for i, scannedRepo := range pCatalog.Status.ScannedRepositories {
		if scannedRepo.URL == repo.URL {
			pCatalog.Status.ScannedRepositories[i].Tags = newTags
			return
		}
	}
	pCatalog.Status.ScannedRepositories = append(pCatalog.Status.ScannedRepositories, profilesv1.ScannedRepository{
		URL:  repo.URL,
		Tags: newTags,
	})
}

func (r *ProfileCatalogSourceReconciler) updateCatalogWithRepositories(ctx context.Context, pCatalog *profilesv1.ProfileCatalogSource, scanner scanner.RepoScanner, logger logr.Logger) error {
	for _, repo := range pCatalog.Spec.Repos {
		logger.Info("scan repo for profiles", "repo", repo)
		var secret *corev1.Secret
		if repo.SecretRef != nil {
			secret = &corev1.Secret{}
			objectKey := client.ObjectKey{Name: repo.SecretRef.Name, Namespace: pCatalog.Namespace}
			if err := r.Client.Get(ctx, objectKey, secret); err != nil {
				return fmt.Errorf("failed to find secret for repo %v: %w", repo, err)
			}
		}
		var alreadyScannedTags []string
		for _, scannedRepo := range pCatalog.Status.ScannedRepositories {
			if scannedRepo.URL == repo.URL {
				alreadyScannedTags = scannedRepo.Tags
			}
		}

		profiles, newTags, err := scanner.ScanRepository(repo, secret, alreadyScannedTags)
		if err != nil {
			return err
		}

		createOrAppendScannedRepositoryStatus(pCatalog, repo, newTags)
		logger.Info("updating catalog with scanning reuslts", "profiles", profiles)
		r.Profiles.Append(pCatalog.Name, profiles...)
	}

	if err := r.Client.Status().Update(ctx, pCatalog); err != nil {
		return fmt.Errorf("failed to patch status: %w", err)
	}
	return nil
}

func createOrAppendScannedRepositoryStatus(pCatalog *profilesv1.ProfileCatalogSource, repo profilesv1.Repository, newTags []string) {
	for i, scannedRepo := range pCatalog.Status.ScannedRepositories {
		if scannedRepo.URL == repo.URL {
			pCatalog.Status.ScannedRepositories[i].Tags = append(scannedRepo.Tags, newTags...)
			return
		}
	}
	pCatalog.Status.ScannedRepositories = append(pCatalog.Status.ScannedRepositories, profilesv1.ScannedRepository{
		URL:  repo.URL,
		Tags: newTags,
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProfileCatalogSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&profilesv1.ProfileCatalogSource{}).
		Complete(r)
}
