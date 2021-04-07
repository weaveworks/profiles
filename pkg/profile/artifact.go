package profile

import (
	"context"
	"fmt"
	"strings"
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
)

type Status struct {
	ResourcesExist     bool
	NotReadyConditions []metav1.Condition
}

// CreateArtifacts creates and inserts objects to the cluster to deploy the
// profile as a HelmRelease.
func (p *Profile) CreateArtifacts(ctx context.Context) error {
	gitRes, helmRes, kustomizeRes, err := p.makeArtifacts()
	if err != nil {
		return err
	}

	p.log.Info("creating GitRepository", "resource", gitRes.ObjectMeta.Name)
	if err := p.client.Create(ctx, gitRes); err != nil {
		return fmt.Errorf("failed to create GitRepository resource: %w", err)
	}

	if helmRes != nil {
		p.log.Info("creating HelmRelease", "resource", helmRes.ObjectMeta.Name)
		if err := p.client.Create(ctx, helmRes); err != nil {
			return fmt.Errorf("failed to create HelmRelease resource: %w", err)
		}
	} else {
		p.log.Info("creating Kustomization", "resource", kustomizeRes.ObjectMeta.Name)
		if err := p.client.Create(ctx, kustomizeRes); err != nil {
			return fmt.Errorf("failed to create Kustomization resource: %w", err)
		}
	}
	p.log.Info("all artifacts created")
	return nil
}

// Checks if the artifacts exists and returns any ready!=true conditions on the artifacts.
func (p *Profile) ArtifactStatus(ctx context.Context) (Status, error) {
	resourcesExist, gitRes, helmRes, kustomizeRes, err := p.getResources(ctx)
	if err != nil {
		return Status{}, err
	}

	return Status{
		ResourcesExist:     resourcesExist,
		NotReadyConditions: p.checkResourcesReady(gitRes, helmRes, kustomizeRes),
	}, nil
}

func (p *Profile) checkResourcesReady(gitRes *sourcev1.GitRepository, helmRes *helmv2.HelmRelease, kustomizeRes *kustomizev1.Kustomization) []metav1.Condition {
	var notReadyConditions []metav1.Condition
	if gitRes != nil {
		if condition := getNotReadyCondition(gitRes.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}

	if helmRes != nil {
		if condition := getNotReadyCondition(helmRes.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}

	if kustomizeRes != nil {
		if condition := getNotReadyCondition(kustomizeRes.Status.Conditions); condition != (metav1.Condition{}) {
			notReadyConditions = append(notReadyConditions, condition)
		}
	}
	return notReadyConditions
}

func (p *Profile) getResources(ctx context.Context) (bool, *sourcev1.GitRepository, *helmv2.HelmRelease, *kustomizev1.Kustomization, error) {
	gitRes, helmRes, kustomizeRes, err := p.makeArtifacts()
	if err != nil {
		return false, nil, nil, nil, err
	}
	gitResExists := true
	helmResExists := true
	kustomizeResExists := true
	if gitRes != nil {
		gitResExists, err = p.getResourceIfExists(ctx, gitRes)
		if err != nil {
			return false, nil, nil, nil, err
		}
	}

	if helmRes != nil {
		helmResExists, err = p.getResourceIfExists(ctx, helmRes)
		if err != nil {
			return false, nil, nil, nil, err
		}
	}

	if kustomizeRes != nil {
		kustomizeResExists, err = p.getResourceIfExists(ctx, kustomizeRes)
		if err != nil {
			return false, nil, nil, nil, err
		}
	}
	return gitResExists && helmResExists && kustomizeResExists, gitRes, helmRes, kustomizeRes, nil
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
	gitRes, helmRes, kustomizeRes, err := p.makeArtifacts()
	if err != nil {
		return nil, err
	}

	objs = append(objs, gitRes)

	if helmRes != nil {
		objs = append(objs, helmRes)
	} else {
		objs = append(objs, kustomizeRes)
	}
	return objs, nil
}

func (p *Profile) makeArtifacts() (*sourcev1.GitRepository, *helmv2.HelmRelease, *kustomizev1.Kustomization, error) {
	gitRes, err := p.makeGitRepository()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create GitRepository resource: %w", err)
	}

	var helmRes *helmv2.HelmRelease
	var kustomizeRes *kustomizev1.Kustomization

	switch kind := p.definition.Spec.Artifacts[0].Kind; kind {
	case "HelmChart":
		helmRes, err = p.makeHelmRelease()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create HelmRelease resource: %w", err)

		}
	case "Kustomize":
		kustomizeRes, err = p.makeKustomization()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create Kustomization resource: %w", err)
		}
	default:
		return nil, nil, nil, fmt.Errorf("artifact kind %q not recognized", kind)
	}

	return gitRes, helmRes, kustomizeRes, nil
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

func (p *Profile) makeHelmRelease() (*helmv2.HelmRelease, error) {
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
					Chart: p.definition.Spec.Artifacts[0].Path,
					SourceRef: helmv2.CrossNamespaceObjectReference{
						Kind:      sourcev1.GitRepositoryKind,
						Name:      p.makeGitRepoName(),
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

func (p *Profile) makeArtifactName() string {
	return join(p.subscription.Name, p.definition.Name, p.definition.Spec.Artifacts[0].Name)
}

func join(s ...string) string {
	return strings.Join(s, "-")
}
