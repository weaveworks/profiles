package profile

import (
	"reflect"
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KustomizeRequiresUpdate checks if the git kustomization resource requires updating
func KustomizeRequiresUpdate(oldRes, newRes *kustomizev1.Kustomization) bool {
	switch {
	case oldRes.Spec.Path != newRes.Spec.Path:
		return true
	case oldRes.Spec.Interval != newRes.Spec.Interval:
		return true
	case oldRes.Spec.Prune != newRes.Spec.Prune:
		return true
	case oldRes.Spec.TargetNamespace != newRes.Spec.TargetNamespace:
		return true
	case !reflect.DeepEqual(oldRes.Spec.SourceRef, newRes.Spec.SourceRef):
		return true
	default:
		return false
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
	return kustomization, nil
}
