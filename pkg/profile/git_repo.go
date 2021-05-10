package profile

import (
	"strings"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// gitRepoRequiresUpdate checks if the git repository resource requires updating
func gitRepoRequiresUpdate(existingRes, newRes *sourcev1.GitRepository) bool {
	switch {
	case existingRes.Spec.URL != newRes.Spec.URL:
		return true
	case existingRes.Spec.Reference.Branch != newRes.Spec.Reference.Branch:
		return true
	default:
		return false
	}
}

func (p *Profile) makeGitRepository() *sourcev1.GitRepository {
	ref := &sourcev1.GitRepositoryRef{
		Branch: p.subscription.Spec.Branch,
	}
	if p.subscription.Spec.Version != "" {
		ref = &sourcev1.GitRepositoryRef{
			Tag: p.subscription.Spec.Version,
		}
	}
	return &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeGitRepoName(),
			Namespace: p.subscription.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.GitRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.GitRepositorySpec{
			URL:       p.subscription.Spec.ProfileURL,
			Reference: ref,
		},
	}
}

func (p *Profile) makeGitRepoName() string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	if p.subscription.Spec.Version != "" {
		parts := strings.Split(p.subscription.Spec.Version, "/")
		return join(p.subscription.Name, repoName, parts[1])
	}
	return join(p.subscription.Name, repoName, p.subscription.Spec.Branch)
}
