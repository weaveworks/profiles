package profile

import (
	"context"
	"fmt"
	"strings"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// location defines where the helm chart is located.
// remote -> in a helm repository
// local -> in the profile repository
type location int

var local location = 0
var remote location = 1

// CreateArtifacts creates and inserts objects to the cluster to deploy the
// profile as a HelmRelease.
func (p *Profile) CreateArtifacts(ctx context.Context) error {
	if err := p.createGitRepository(ctx); err != nil {
		return fmt.Errorf("failed to create GitRepository resource: %w", err)
	}

	// TODO don't depend on 0.
	switch kind := p.definition.Spec.Artifacts[0].Kind; kind {
	case profilesv1.HelmChartLocalKind:
		if err := p.createHelmRelease(ctx, local); err != nil {
			return fmt.Errorf("failed to create HelmRelease resource: %w", err)
		}
	case profilesv1.HelmChartRemoteKind:
		if p.definition.Spec.Artifacts[0].HelmURL == "" {
			return fmt.Errorf("helmURL for helm repository for remote helm artifact must not be empty")
		}
		if err := p.createHelmRepository(ctx, p.definition.Spec.Artifacts[0].HelmURL); err != nil {
			return fmt.Errorf("failed to create HelmRepository resource: %w", err)
		}
		if err := p.createHelmRelease(ctx, remote); err != nil {
			return fmt.Errorf("failed to create HelmRelease for remote resource: %w", err)
		}
	case profilesv1.KustomizeKind:
		if err := p.createKustomization(ctx); err != nil {
			return fmt.Errorf("failed to create Kustomization resource: %w", err)
		}
	default:
		return fmt.Errorf("artifact kind %q not recognized", kind)
	}

	p.log.Info("all artifacts created")
	return nil
}

// MakeArtifacts creates and returns a slice of runtime.Object values, which if
// applied to a cluster would deploy the profile as a HelmRelease.
func (p *Profile) MakeArtifacts() ([]runtime.Object, error) {
	objs := []runtime.Object{}
	gr, err := p.makeGitRepository()
	if err != nil {
		return nil, err
	}

	hr, err := p.makeHelmRelease(local)
	if err != nil {
		return nil, err
	}

	objs = append(objs, gr, hr)
	return objs, nil
}

func (p *Profile) makeGitRepository() (*sourcev1.GitRepository, error) {
	gitRepo := &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeRepoName(),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.GitRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: p.subscription.Spec.ProfileURL,
			Reference: &sourcev1.GitRepositoryRef{
				Branch: p.subscription.Spec.Branch,
			},
		},
	}
	err := controllerutil.SetControllerReference(&p.subscription, gitRepo, p.client.Scheme())
	if err != nil {
		return nil, fmt.Errorf("failed to set resource ownership on %s: %w", gitRepo.ObjectMeta.Name, err)
	}
	return gitRepo, nil
}

func (p *Profile) makeHelmRepository(url string) (*sourcev1.HelmRepository, error) {
	helmRepo := &sourcev1.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeRepoName(),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.HelmRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.HelmRepositorySpec{
			URL: url,
		},
	}
	err := controllerutil.SetControllerReference(&p.subscription, helmRepo, p.client.Scheme())
	if err != nil {
		return nil, fmt.Errorf("failed to set resource ownership on %s: %w", helmRepo.ObjectMeta.Name, err)
	}
	return helmRepo, nil
}

func (p *Profile) createGitRepository(ctx context.Context) error {
	r, err := p.makeGitRepository()
	if err != nil {
		return err
	}
	p.log.Info("creating GitRepository", "resource", r.ObjectMeta.Name)
	return p.client.Create(ctx, r)
}

func (p *Profile) createHelmRepository(ctx context.Context, url string) error {
	r, err := p.makeHelmRepository(url)
	if err != nil {
		return err
	}
	p.log.Info("creating HelmRepository", "resource", r.ObjectMeta.Name)
	return p.client.Create(ctx, r)
}

func (p *Profile) createHelmRelease(ctx context.Context, loc location) error {
	r, err := p.makeHelmRelease(loc)
	if err != nil {
		return err
	}
	p.log.Info("creating remote HelmRelease", "resource", r.ObjectMeta.Name)
	return p.client.Create(ctx, r)
}

func (p *Profile) makeHelmRelease(loc location) (*helmv2.HelmRelease, error) {
	kind := sourcev1.GitRepositoryKind
	chart := p.definition.Spec.Artifacts[0].Path
	if loc == remote {
		kind = sourcev1.HelmRepositoryKind
		chart = p.definition.Spec.Artifacts[0].Name
	}
	helmRelease := &helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeArtifactName(),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       helmv2.HelmReleaseKind,
			APIVersion: helmv2.GroupVersion.String(),
		},
		Spec: helmv2.HelmReleaseSpec{
			Chart: helmv2.HelmChartTemplate{
				Spec: helmv2.HelmChartTemplateSpec{
					// TODO obvs don't rely on index 0
					Chart: chart,
					SourceRef: helmv2.CrossNamespaceObjectReference{
						Kind:      kind,
						Name:      p.makeRepoName(),
						Namespace: p.subscription.ObjectMeta.Namespace,
					},
				},
			},
			Values:     p.subscription.Spec.Values,
			ValuesFrom: p.subscription.Spec.ValuesFrom,
		},
	}
	err := controllerutil.SetControllerReference(&p.subscription, helmRelease, p.client.Scheme())
	if err != nil {
		return nil, fmt.Errorf("failed to set resource ownership: %w", err)
	}
	return helmRelease, nil
}

func (p *Profile) makeKustomization() (*kustomizev1.Kustomization, error) {
	kustomization := &kustomizev1.Kustomization{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeArtifactName(),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       kustomizev1.KustomizationKind,
			APIVersion: kustomizev1.GroupVersion.String(),
		},
		Spec: kustomizev1.KustomizationSpec{
			// TODO obvs don't rely on index 0
			Path:            p.definition.Spec.Artifacts[0].Path,
			Interval:        metav1.Duration{Duration: time.Minute * 5},
			Prune:           true,
			TargetNamespace: p.subscription.ObjectMeta.Namespace,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      sourcev1.GitRepositoryKind,
				Name:      p.makeRepoName(),
				Namespace: p.subscription.ObjectMeta.Namespace,
			},
		},
	}
	err := controllerutil.SetControllerReference(&p.subscription, kustomization, p.client.Scheme())
	if err != nil {
		return nil, fmt.Errorf("failed to set resource ownership: %w", err)
	}
	return kustomization, nil
}

func (p *Profile) createKustomization(ctx context.Context) error {
	r, err := p.makeKustomization()
	if err != nil {
		return err
	}
	p.log.Info("creating Kustomization", "resource", r.ObjectMeta.Name)
	return p.client.Create(ctx, r)
}

func (p *Profile) makeRepoName() string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(p.subscription.Name, repoName, p.subscription.Spec.Branch)
}

func (p *Profile) makeArtifactName() string {
	return join(p.subscription.Name, p.definition.Name, p.definition.Spec.Artifacts[0].Name)
}

func join(s ...string) string {
	return strings.Join(s, "-")
}
