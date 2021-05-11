package profile

import (
	"path"
	"reflect"
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// kustomizeRequiresUpdate checks if the git kustomization resource requires updating
func kustomizeRequiresUpdate(existingRes, desiredRes *kustomizev1.Kustomization) bool {
	switch {
	case existingRes.Spec.Path != desiredRes.Spec.Path:
		return true
	case existingRes.Spec.Interval != desiredRes.Spec.Interval:
		return true
	case existingRes.Spec.Prune != desiredRes.Spec.Prune:
		return true
	case existingRes.Spec.TargetNamespace != desiredRes.Spec.TargetNamespace:
		return true
	case !reflect.DeepEqual(existingRes.Spec.SourceRef, desiredRes.Spec.SourceRef):
		return true
	default:
		return false
	}
}

func (p *Profile) makeKustomization(artifact profilesv1.Artifact, repoPath string) *kustomizev1.Kustomization {
	return &kustomizev1.Kustomization{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeArtifactName(artifact.Name),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       kustomizev1.KustomizationKind,
			APIVersion: kustomizev1.GroupVersion.String(),
		},
		Spec: kustomizev1.KustomizationSpec{
			Path:            path.Join(repoPath, artifact.Path),
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
}
