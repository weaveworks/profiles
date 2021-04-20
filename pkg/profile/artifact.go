package profile

import (
	"context"
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
func (p *Profile) ReconcileArtifacts(ctx context.Context) error {
	objs, err := p.MakeArtifacts()
	if err != nil {
		return err
	}

	for _, o := range objs {
		obj, ok := o.(client.Object)
		if !ok {
			return fmt.Errorf("object %v cannot be asserted to client.Object", o)
		}
		if err := p.reconcileArtifact(ctx, obj); err != nil {
			return err
		}
	}

	p.log.Info("all artifacts created")
	return nil
}
func (p *Profile) reconcileArtifact(ctx context.Context, obj client.Object) error {
	oldObj := obj.DeepCopyObject().(client.Object)
	err := p.client.Get(ctx, client.ObjectKeyFromObject(obj), oldObj)
	if err == nil {
		return p.updateResource(ctx, oldObj, obj)
	} else if apierrors.IsNotFound(err) {
		p.log.Info("creating...", "kind", obj.GetObjectKind().GroupVersionKind().Kind, "resource", obj.GetName())
		if err := p.client.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create %s: %w", obj.GetObjectKind().GroupVersionKind().Kind, err)
		}
		return nil
	} else {
		//todo update
		return fmt.Errorf("failed to create %s: %w", obj.GetObjectKind().GroupVersionKind().Kind, err)
	}
}

func (p *Profile) updateResource(ctx context.Context, oldRes, newRes client.Object) error {
	switch newRes := newRes.(type) {
	case *sourcev1.GitRepository:
		if !GitRepoRequiresUpdate(oldRes.(*sourcev1.GitRepository), newRes) {
			return nil
		}
		oldRes.(*sourcev1.GitRepository).Spec = newRes.Spec
	case *sourcev1.HelmRepository:
		if !HelmRepoRequiresUpdate(oldRes.(*sourcev1.HelmRepository), newRes) {
			return nil
		}
		oldRes.(*sourcev1.HelmRepository).Spec = newRes.Spec
	case *helmv2.HelmRelease:
		if !HelmReleaseRequiresUpdate(oldRes.(*helmv2.HelmRelease), newRes) {
			return nil
		}
		oldRes.(*helmv2.HelmRelease).Spec = newRes.Spec
	case *kustomizev1.Kustomization:
		if !KustomizeRequiresUpdate(oldRes.(*kustomizev1.Kustomization), newRes) {
			return nil
		}
		oldRes.(*kustomizev1.Kustomization).Spec = newRes.Spec
	default:
		return nil
	}
	p.log.Info(fmt.Sprintf("updating %s, %s", oldRes.GetObjectKind(), oldRes.GetName()))
	return p.client.Update(ctx, oldRes)
}

// ArtifactStatus checks if the artifacts exists and returns any ready!=true conditions on the artifacts.
func (p *Profile) ArtifactStatus(ctx context.Context) (Status, error) {
	resourcesExist, objs, err := p.getResources(ctx)
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

func (p *Profile) getResources(ctx context.Context) (bool, []runtime.Object, error) {
	objs, err := p.MakeArtifacts()
	if err != nil {
		return false, nil, err
	}
	for _, o := range objs {
		obj, ok := o.(client.Object)
		if !ok {
			return false, nil, fmt.Errorf("object is not a client.Object %v", o)
		}
		if exists, err := p.getResourceIfExists(ctx, obj); !exists || err != nil {
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

func (p *Profile) getResourceIfExists(ctx context.Context, res client.Object) (bool, error) {
	if err := p.client.Get(ctx, client.ObjectKey{Name: res.GetName(), Namespace: res.GetNamespace()}, res); err != nil {
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
	var (
		objs   []runtime.Object
		gitRes *sourcev1.GitRepository
	)

	for _, artifact := range p.definition.Spec.Artifacts {
		if err := artifact.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed for artifact %s: %w", artifact.Name, err)
		}
		switch artifact.Kind {
		case profilesv1.HelmChartKind:
			helmRes, err := p.makeHelmRelease(artifact)
			if err != nil {
				return nil, fmt.Errorf("failed to create HelmRelease resource: %w", err)
			}
			objs = append(objs, helmRes)
			if artifact.Path != "" && gitRes == nil {
				// this resource is added at the end because it's generated once.
				gitRes, err = p.makeGitRepository()
				if err != nil {
					return nil, fmt.Errorf("failed to create GitRepository resource: %w", err)
				}
			}
			if artifact.Chart != nil {
				helmRep, err := p.makeHelmRepository(artifact.Chart.URL, artifact.Chart.Name)
				if err != nil {
					return nil, fmt.Errorf("failed to create HelmRepository resource: %w", err)
				}
				objs = append(objs, helmRep)
			}
		case profilesv1.KustomizeKind:
			kustomizeRes, err := p.makeKustomization(artifact)
			if err != nil {
				return nil, fmt.Errorf("failed to create Kustomization resource: %w", err)
			}
			objs = append(objs, kustomizeRes)

			if gitRes == nil {
				// this resource is added at the end because it's generated once.
				gitRes, err = p.makeGitRepository()
				if err != nil {
					return nil, fmt.Errorf("failed to create GitRepository resource: %w", err)
				}
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
