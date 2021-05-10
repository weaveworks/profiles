package profile

import (
	"fmt"
	"strings"

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

// ReconcileArtifacts ensures the artifact resources are applied to the cluster
func (p *Profile) ReconcileArtifacts() error {
	objs, err := p.MakeArtifacts()
	if err != nil {
		return err
	}

	for _, o := range objs {
		obj, ok := o.(client.Object)
		if !ok {
			return fmt.Errorf("object %v cannot be asserted to client.Object", o)
		}
		if err := p.reconcileArtifact(obj); err != nil {
			return err
		}
	}

	p.log.Info("all artifacts created")
	return nil
}
func (p *Profile) reconcileArtifact(desiredObj client.Object) error {
	existingObj := desiredObj.DeepCopyObject().(client.Object)
	err := p.client.Get(p.ctx, client.ObjectKeyFromObject(desiredObj), existingObj)
	if err == nil {
		return p.updateResource(existingObj, desiredObj)
	}
	if apierrors.IsNotFound(err) {
		p.log.Info("creating...", "kind", desiredObj.GetObjectKind().GroupVersionKind().Kind, "resource", desiredObj.GetName())
		if err := p.client.Create(p.ctx, desiredObj); err != nil {
			return fmt.Errorf("failed to create %s: %w", desiredObj.GetObjectKind().GroupVersionKind().Kind, err)
		}
		return nil
	}
	return fmt.Errorf("failed to get resource %s %s/%s: %w", desiredObj.GetObjectKind().GroupVersionKind().Kind, desiredObj.GetNamespace(), desiredObj.GetName(), err)
}

func (p *Profile) updateResource(existingObj, desiredObj client.Object) error {
	switch desiredObj := desiredObj.(type) {
	case *sourcev1.GitRepository:
		if !gitRepoRequiresUpdate(existingObj.(*sourcev1.GitRepository), desiredObj) {
			return nil
		}
		existingObj.(*sourcev1.GitRepository).Spec = desiredObj.Spec
	case *sourcev1.HelmRepository:
		if !helmRepoRequiresUpdate(existingObj.(*sourcev1.HelmRepository), desiredObj) {
			return nil
		}
		existingObj.(*sourcev1.HelmRepository).Spec = desiredObj.Spec
	case *helmv2.HelmRelease:
		if !helmReleaseRequiresUpdate(existingObj.(*helmv2.HelmRelease), desiredObj) {
			return nil
		}
		existingObj.(*helmv2.HelmRelease).Spec = desiredObj.Spec
	case *kustomizev1.Kustomization:
		if !kustomizeRequiresUpdate(existingObj.(*kustomizev1.Kustomization), desiredObj) {
			return nil
		}
		existingObj.(*kustomizev1.Kustomization).Spec = desiredObj.Spec
	default:
		return nil
	}
	p.log.Info(fmt.Sprintf("updating %s, %s", existingObj.GetObjectKind(), existingObj.GetName()))
	return p.client.Update(p.ctx, existingObj)
}

// ArtifactStatus checks if the artifacts exists and returns any ready!=true conditions on the artifacts.
func (p *Profile) ArtifactStatus() (Status, error) {
	resourcesExist, objs, err := p.getResources()
	if err != nil {
		return Status{}, err
	}

	conditions, err := p.checkResourcesReady(objs)
	if err != nil {
		return Status{}, err
	}
	return Status{
		ResourcesExist:     resourcesExist,
		NotReadyConditions: conditions,
	}, nil
}

func (p *Profile) checkResourcesReady(objs []runtime.Object) ([]metav1.Condition, error) {
	var notReadyConditions []metav1.Condition
	for _, o := range objs {
		switch t := o.(type) {
		// annoying, but can't combine these.
		case *sourcev1.GitRepository:
			if condition := getNotReadyCondition(t.Status.Conditions); condition != (metav1.Condition{}) {
				notReadyConditions = append(notReadyConditions, condition)
			}
		case *sourcev1.HelmRepository:
			if condition := getNotReadyCondition(t.Status.Conditions); condition != (metav1.Condition{}) {
				notReadyConditions = append(notReadyConditions, condition)
			}
		case *helmv2.HelmRelease:
			if condition := getNotReadyCondition(t.Status.Conditions); condition != (metav1.Condition{}) {
				notReadyConditions = append(notReadyConditions, condition)
			}
		case *kustomizev1.Kustomization:
			if condition := getNotReadyCondition(t.Status.Conditions); condition != (metav1.Condition{}) {
				notReadyConditions = append(notReadyConditions, condition)
			}
		default:
			return nil, fmt.Errorf("unsupported resource type %v", t)
		}

	}
	return notReadyConditions, nil
}

func (p *Profile) getResources() (bool, []runtime.Object, error) {
	objs, err := p.MakeArtifacts()
	if err != nil {
		return false, nil, err
	}
	for _, o := range objs {
		obj, ok := o.(client.Object)
		if !ok {
			return false, nil, fmt.Errorf("object is not a client.Object %v", o)
		}
		if exists, err := p.getResourceIfExists(obj); !exists || err != nil {
			return false, nil, err
		}
	}
	return true, objs, nil
}

func getNotReadyCondition(conditions []metav1.Condition) metav1.Condition {
	for _, condition := range conditions {
		if condition.Type == "Ready" && string(condition.Status) != "True" {
			return condition
		}
	}
	return metav1.Condition{}
}

func (p *Profile) getResourceIfExists(res client.Object) (bool, error) {
	if err := p.client.Get(p.ctx, client.ObjectKey{Name: res.GetName(), Namespace: res.GetNamespace()}, res); err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get resource: %w", err)
	}

	return true, nil
}

// MakeOwnerlessArtifacts generates artifacts without owners for manual applying to
// a personal cluster.
func (p *Profile) MakeOwnerlessArtifacts() ([]runtime.Object, error) {
	return p.makeOwnerlessArtifacts([]string{p.profileRepo()})
}

func (p *Profile) profileRepo() string {
	return p.subscription.Spec.ProfileURL + ":" + p.subscription.Spec.Branch
}

func (p *Profile) makeOwnerlessArtifacts(profileRepos []string) ([]runtime.Object, error) {
	var (
		objs   []runtime.Object
		gitRes *sourcev1.GitRepository
	)

	for _, artifact := range p.definition.Spec.Artifacts {
		if err := artifact.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed for artifact %s: %w", artifact.Name, err)
		}
		switch artifact.Kind {
		case profilesv1.ProfileKind:
			nestedProfileDef, err := getProfileDefinition(artifact.Profile.URL, artifact.Profile.Branch, p.log)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch profile %q: %w", artifact.Name, err)
			}
			nestedProfile := p.subscription.DeepCopyObject().(*profilesv1.ProfileSubscription)
			nestedProfile.Spec.ProfileURL = artifact.Profile.URL
			nestedProfile.Spec.Branch = artifact.Profile.Branch

			nestedSub := New(p.ctx, nestedProfileDef, *nestedProfile, p.client, p.log)
			profileRepoName := nestedSub.profileRepo()
			if containsKey(profileRepos, profileRepoName) {
				return nil, fmt.Errorf("recursive artifact detected: profile %s on branch %s contains an artifact that points recursively back at itself", artifact.Profile.URL, artifact.Profile.Branch)
			}
			profileRepos = append(profileRepos, profileRepoName)
			nestedObjs, err := nestedSub.makeOwnerlessArtifacts(profileRepos)
			if err != nil {
				return nil, fmt.Errorf("failed to generate resources for nested profile %q: %w", artifact.Name, err)
			}
			objs = append(objs, nestedObjs...)
		case profilesv1.HelmChartKind:
			objs = append(objs, p.makeHelmRelease(artifact))
			if artifact.Path != "" && gitRes == nil {
				// this resource is added at the end because it's generated once.
				gitRes = p.makeGitRepository()
			}
			if artifact.Chart != nil {
				objs = append(objs, p.makeHelmRepository(artifact.Chart.URL, artifact.Chart.Name))
			}
		case profilesv1.KustomizeKind:
			objs = append(objs, p.makeKustomization(artifact))
			if gitRes == nil {
				// this resource is added at the end because it's generated once.
				gitRes = p.makeGitRepository()
			}
		default:
			return nil, fmt.Errorf("artifact kind %q not recognized", artifact.Kind)
		}
	}

	// Add the git res as the first object to be created.
	if gitRes != nil {
		objs = append([]runtime.Object{gitRes}, objs...)
	}
	return objs, nil
}

func containsKey(list []string, key string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}
	return false
}

// MakeArtifacts creates and returns a slice of runtime.Object values, which if
// applied to a cluster would deploy the profile as a HelmRelease.
func (p *Profile) MakeArtifacts() ([]runtime.Object, error) {
	objs, err := p.MakeOwnerlessArtifacts()
	if err != nil {
		return nil, err
	}

	// setup ownership
	for _, o := range objs {
		obj, ok := o.(client.Object)
		if !ok {
			return nil, fmt.Errorf("object is not a client.Object %v", o)
		}
		if err := controllerutil.SetControllerReference(&p.subscription, obj, p.client.Scheme()); err != nil {
			return nil, fmt.Errorf("failed to set resource ownership on %s: %w", obj.GetName(), err)
		}
	}

	return objs, nil
}

func (p *Profile) makeArtifactName(name string) string {
	return join(p.subscription.Name, p.definition.Name, name)
}

func join(s ...string) string {
	return strings.Join(s, "-")
}
