package profile

import (
	"context"
	"fmt"
	"strings"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	gitRepositoryKind       = "GitRepository"
	gitRepositoryAPIVersion = "source.toolkit.fluxcd.io/v1beta1"
	helmReleaseKind         = "HelmRelease"
	helmReleaseAPIVersion   = "helm.toolkit.fluxcd.io/v2beta1"
)

func (p *Profile) CreateArtifacts(ctx context.Context) error {
	if err := p.createGitRepository(ctx); err != nil {
		return errors.Wrap(err, "failed to create GitRepository resource")
	}

	if err := p.createHelmRelease(ctx); err != nil {
		return errors.Wrap(err, "failed to create HelmRelease resource")
	}

	p.log.Info("all artifacts created")
	return nil
}

func (p *Profile) createGitRepository(ctx context.Context) error {
	gitRefName := p.makeGitRepoName()
	namespace := p.subscription.Namespace
	gitRepo := sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gitRefName,
			Namespace: namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       gitRepositoryKind,
			APIVersion: gitRepositoryAPIVersion,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: p.subscription.Spec.ProfileURL,
			Reference: &sourcev1.GitRepositoryRef{
				Branch: p.subscription.Spec.Branch,
			},
		},
	}

	p.log.Info(fmt.Sprintf("creating git repository %s/%s", namespace, gitRefName))
	return p.client.Create(ctx, &gitRepo)
}

func (p *Profile) createHelmRelease(ctx context.Context) error {
	namespace := p.subscription.Namespace
	helmReleasename := p.makeHelmReleaseName()
	helmRelease := helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmReleasename,
			Namespace: namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       helmReleaseKind,
			APIVersion: helmReleaseAPIVersion,
		},
		Spec: helmv2.HelmReleaseSpec{
			Chart: helmv2.HelmChartTemplate{
				Spec: helmv2.HelmChartTemplateSpec{
					// TODO obvs don't rely on index 0
					Chart: p.definition.Spec.Artifacts[0].Path,
					SourceRef: helmv2.CrossNamespaceObjectReference{
						Kind:      gitRepositoryKind,
						Name:      p.makeGitRepoName(),
						Namespace: namespace,
					},
				},
			},
		},
	}

	p.log.Info(fmt.Sprintf("creating helm release %s/%s", namespace, helmReleasename))
	return p.client.Create(ctx, &helmRelease)
}

func (p *Profile) makeGitRepoName() string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return fmt.Sprintf("%s-%s-%s", p.subscription.Name, repoName, p.subscription.Spec.Branch)
}

func (p *Profile) makeHelmReleaseName() string {
	return fmt.Sprintf("%s-%s-%s", p.subscription.Name, p.definition.Name, p.definition.Spec.Artifacts[0].Name)
}
