package profile

import (
	"strings"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitRepoRequiresUpdate checks if the git repository resource requires updating
func GitRepoRequiresUpdate(oldRes, newRes *sourcev1.GitRepository) bool {
	switch {
	case oldRes.Spec.URL != newRes.Spec.URL:
		return true
	case oldRes.Spec.Reference.Branch != newRes.Spec.Reference.Branch:
		return true
	default:
		return false
	}
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
	return gitRepo, nil
}

func (p *Profile) makeGitRepoName() string {
	repoParts := strings.Split(p.subscription.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(p.subscription.Name, repoName, p.subscription.Spec.Branch)
}
