package profile

import (
	"context"
	"fmt"
	"strings"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// Status contains the status of the artifacts
type Status struct {
	ResourcesExist     bool
	NotReadyConditions []metav1.Condition
}

// CreateArtifacts generate and creates the objects in the cluster to deploy the
// profile
func (p *Profile) CreateArtifacts(ctx context.Context) error {
	gitRes, helmRepositories, helmResources, kustomizeResources, err := p.makeArtifacts()
	if err != nil {
		return err
	}

	p.log.Info("creating GitRepository", "resource", gitRes.ObjectMeta.Name)
	if err := p.client.Create(ctx, gitRes); err != nil {
		return fmt.Errorf("failed to create GitRepository resource: %w", err)
	}

	for _, helmRep := range helmRepositories {
		p.log.Info("creating HelmRepositories", "resource", helmRep.ObjectMeta.Name)
		if err := p.client.Create(ctx, helmRep); err != nil {
			return fmt.Errorf("failed to create HelmRepository resource: %w", err)
		}
	}
	for _, helmRes := range helmResources {
		p.log.Info("creating HelmRelease", "resource", helmRes.ObjectMeta.Name)
		if err := p.client.Create(ctx, helmRes); err != nil {
			return fmt.Errorf("failed to create HelmRelease resource: %w", err)
		}
	}

	for _, kustomizeRes := range kustomizeResources {
		p.log.Info("creating Kustomization", "resource", kustomizeRes.ObjectMeta.Name)
		if err := p.client.Create(ctx, kustomizeRes); err != nil {
			return fmt.Errorf("failed to create Kustomization resource: %w", err)
		}
	}
	p.log.Info("all artifacts created")
	return nil
}

// ArtifactStatus checks if the artifacts exists and returns any ready!=true conditions on the artifacts.
func (p *Profile) ArtifactStatus(ctx context.Context) (Status, error) {
	resourcesExist, gitRes, helmRepositories, helmResources, kustomizeResources, err := p.getResources(ctx)
	if err != nil {
		return Status{}, err
	}

	return Status{
		ResourcesExist:     resourcesExist,
		NotReadyConditions: p.checkResourcesReady(gitRes, helmRepositories, helmResources, kustomizeResources),
	}, nil
}

func (p *Profile) checkResourcesReady(gitRes *sourcev1.GitRepository, helmRepositories []*sourcev1.HelmRepository, helmResources []*helmv2.HelmRelease, kustomizeResources []*kustomizev1.Kustomization) []metav1.Condition {
	var notReadyConditions []metav1.Condition
	if gitRes != nil {
		if condition := getNotReadyCondition(gitRes.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}

	for _, helmRep := range helmRepositories {
		if condition := getNotReadyCondition(helmRep.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}

	for _, helmRes := range helmResources {
		if condition := getNotReadyCondition(helmRes.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}

	for _, kustomizeRes := range kustomizeResources {
		if condition := getNotReadyCondition(kustomizeRes.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}
	return notReadyConditions
}

func (p *Profile) getResources(ctx context.Context) (bool, *sourcev1.GitRepository, []*sourcev1.HelmRepository, []*helmv2.HelmRelease, []*kustomizev1.Kustomization, error) {
	gitRes, helmRepositories, helmResources, kustomizeResources, err := p.makeArtifacts()
	if err != nil {
		return false, nil, nil, nil, nil, err
	}
	if gitRes != nil {
		gitResExists, err := p.getResourceIfExists(ctx, gitRes)
		if err != nil || !gitResExists {
			return false, nil, nil, nil, nil, err
		}
	}

	for _, helmRep := range helmRepositories {
		helmRepExists, err := p.getResourceIfExists(ctx, helmRep)
		if err != nil || !helmRepExists {
			return false, nil, nil, nil, nil, err
		}
	}
	for _, helmRes := range helmResources {
		helmResExists, err := p.getResourceIfExists(ctx, helmRes)
		if err != nil || !helmResExists {
			return false, nil, nil, nil, nil, err
		}
	}

	for _, kustomizeRes := range kustomizeResources {
		kustomizeResExists, err := p.getResourceIfExists(ctx, kustomizeRes)
		if err != nil || !kustomizeResExists {
			return false, nil, nil, nil, nil, err
		}
	}
	return true, gitRes, helmRepositories, helmResources, kustomizeResources, nil
}

func getNotReadyCondition(conditions []metav1.Condition) metav1.Condition {
	for _, condition := range conditions {
		if condition.Type == "Ready" && string(condition.Status) != "True" {
			return condition
		}
	}
	return metav1.Condition{}
}

func (p *Profile) getResourceIfExists(ctx context.Context, res client.Object) (bool, error) {
	if err := p.client.Get(ctx, client.ObjectKey{Name: res.GetName(), Namespace: res.GetNamespace()}, res); err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get resource: %w", err)
	}

	return true, nil
}

// MakeArtifacts creates and returns a slice of runtime.Object values, which if
// applied to a cluster would deploy the profile as a HelmRelease.
func (p *Profile) MakeArtifacts() ([]runtime.Object, error) {
	objs := []runtime.Object{}
	gitRes, helmRepositories, helmResources, kustomizeResources, err := p.makeArtifacts()
	if err != nil {
		return nil, err
	}

	objs = append(objs, gitRes)
	for _, helmRepo := range helmRepositories {
		objs = append(objs, helmRepo)
	}
	for _, helmRes := range helmResources {
		objs = append(objs, helmRes)
	}
	for _, kustomizeRes := range kustomizeResources {
		objs = append(objs, kustomizeRes)
	}
	return objs, nil
}

func (p *Profile) makeArtifacts() (*sourcev1.GitRepository, []*sourcev1.HelmRepository, []*helmv2.HelmRelease, []*kustomizev1.Kustomization, error) {
	var helmResources []*helmv2.HelmRelease
	var helmReps []*sourcev1.HelmRepository
	var kustomizeResources []*kustomizev1.Kustomization
	var gitRes *sourcev1.GitRepository

	for _, artifact := range p.definition.Spec.Artifacts {
		switch artifact.Kind {
		case profilesv1.HelmChartKind:
			helmRes, err := p.makeHelmRelease(artifact)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to create HelmRelease resource: %w", err)

			}
			helmResources = append(helmResources, helmRes)
			if artifact.Path != "" {
				gitRes, err = p.makeGitRepository()
				if err != nil {
					return nil, nil, nil, nil, fmt.Errorf("failed to create GitRepository resource: %w", err)
				}
			} else if artifact.Chart != nil {
				helmRep, err := p.makeHelmRepository(artifact.Chart.HelmURL, artifact.Chart.HelmChart)
				if err != nil {
					return nil, nil, nil, nil, fmt.Errorf("failed to create HelmRepository resource: %w", err)
				}
				helmReps = append(helmReps, helmRep)
			}
		case profilesv1.KustomizeKind:
			kustomizeRes, err := p.makeKustomization(artifact)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to create Kustomization resource: %w", err)
			}
			kustomizeResources = append(kustomizeResources, kustomizeRes)

			gitRes, err = p.makeGitRepository()
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to create GitRepository resource: %w", err)
			}
		default:
			return nil, nil, nil, nil, fmt.Errorf("artifact kind %q not recognized", artifact.Kind)
		}
	}

	return gitRes, helmReps, helmResources, kustomizeResources, nil
}

func (p *Profile) makeGitRepository() (*sourcev1.GitRepository, error) {
	gitRepo := &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeGitRepoName(),
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

func (p *Profile) makeHelmRepository(url string, name string) (*sourcev1.HelmRepository, error) {
	helmRepo := &sourcev1.HelmRepository{
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
	err := controllerutil.SetControllerReference(&p.subscription, helmRepo, p.client.Scheme())
	if err != nil {
		return nil, fmt.Errorf("failed to set resource ownership on %s: %w", helmRepo.ObjectMeta.Name, err)
	}
	return helmRepo, nil
}

func (p *Profile) makeHelmRepoName(name string) string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(p.subscription.Name, repoName, p.subscription.Spec.Branch, name)
}

func (p *Profile) makeHelmRelease(artifact profilesv1.Artifact) (*helmv2.HelmRelease, error) {
	var helmChartSpec helmv2.HelmChartTemplateSpec
	if artifact.Path != "" {
		helmChartSpec = p.makeGitChartSpec(artifact.Path)
	} else if artifact.Chart != nil {
		helmChartSpec = p.makeHelmChartSpec(artifact.Chart.HelmChart, artifact.Chart.HelmChartVersion)
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
	err := controllerutil.SetControllerReference(&p.subscription, helmRelease, p.client.Scheme())
	if err != nil {
		return nil, fmt.Errorf("failed to set resource ownership: %w", err)
	}
	return helmRelease, nil
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

func (p *Profile) makeKustomization(artifact profilesv1.Artifact) (*kustomizev1.Kustomization, error) {
	kustomization := &kustomizev1.Kustomization{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeArtifactName(artifact.Name),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       kustomizev1.KustomizationKind,
			APIVersion: kustomizev1.GroupVersion.String(),
		},
		Spec: kustomizev1.KustomizationSpec{
			Path:            artifact.Path,
			Interval:        metav1.Duration{Duration: time.Minute * 5},
			Prune:           true,
			TargetNamespace: p.subscription.ObjectMeta.Namespace,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      sourcev1.GitRepositoryKind,
				Name:      p.makeGitRepoName(),
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

func (p *Profile) makeGitRepoName() string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(p.subscription.Name, repoName, p.subscription.Spec.Branch)
}

func (p *Profile) makeArtifactName(name string) string {
	return join(p.subscription.Name, p.definition.Name, name)
}

func join(s ...string) string {
	return strings.Join(s, "-")
}
