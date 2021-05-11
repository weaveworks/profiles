package profile

import (
	"path"
	"reflect"
	"strings"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// helmReleaseRequiresUpdate checks if the helm release resource requires updating
func helmReleaseRequiresUpdate(existingRes, newRes *helmv2.HelmRelease) bool {
	switch {
	case existingRes.Spec.Chart.Spec.Chart != newRes.Spec.Chart.Spec.Chart:
		return true
	case existingRes.Spec.Chart.Spec.Version != newRes.Spec.Chart.Spec.Version:
		return true
	case !reflect.DeepEqual(existingRes.Spec.Chart.Spec.SourceRef, newRes.Spec.Chart.Spec.SourceRef):
		return true
	case !reflect.DeepEqual(existingRes.Spec.Values, newRes.Spec.Values):
		return true
	case !reflect.DeepEqual(existingRes.Spec.ValuesFrom, newRes.Spec.ValuesFrom):
		return true
	default:
		return false
	}
}

// HelmRepoRequiresUpdate checks if the helm repository resource requires updating
func helmRepoRequiresUpdate(existingRes, newRes *sourcev1.HelmRepository) bool {
	switch {
	case existingRes.Spec.URL != newRes.Spec.URL:
		return true
	default:
		return false
	}
}

func (p *Profile) makeHelmRepository(url string, name string) *sourcev1.HelmRepository {
	return &sourcev1.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeHelmRepoName(name),
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
}

func (p *Profile) makeHelmRepoName(name string) string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(p.subscription.Name, repoName, name)
}

func (p *Profile) makeHelmRelease(artifact profilesv1.Artifact, repoPath string) *helmv2.HelmRelease {
	var helmChartSpec helmv2.HelmChartTemplateSpec
	if artifact.Path != "" {
		helmChartSpec = p.makeGitChartSpec(path.Join(repoPath, artifact.Path))
	} else if artifact.Chart != nil {
		helmChartSpec = p.makeHelmChartSpec(artifact.Chart.Name, artifact.Chart.Version)
	}
	helmRelease := &helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeArtifactName(artifact.Name),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       helmv2.HelmReleaseKind,
			APIVersion: helmv2.GroupVersion.String(),
		},
		Spec: helmv2.HelmReleaseSpec{
			Chart: helmv2.HelmChartTemplate{
				Spec: helmChartSpec,
			},
			Values:     p.subscription.Spec.Values,
			ValuesFrom: p.subscription.Spec.ValuesFrom,
		},
	}
	return helmRelease
}

func (p *Profile) makeGitChartSpec(path string) helmv2.HelmChartTemplateSpec {
	return helmv2.HelmChartTemplateSpec{
		Chart: path,
		SourceRef: helmv2.CrossNamespaceObjectReference{
			Kind:      sourcev1.GitRepositoryKind,
			Name:      p.makeGitRepoName(),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
	}
}

func (p *Profile) makeHelmChartSpec(chart string, version string) helmv2.HelmChartTemplateSpec {
	return helmv2.HelmChartTemplateSpec{
		Chart: chart,
		SourceRef: helmv2.CrossNamespaceObjectReference{
			Kind:      sourcev1.HelmRepositoryKind,
			Name:      p.makeHelmRepoName(chart),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		Version: version,
	}
}
