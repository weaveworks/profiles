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
	return &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.makeGitRepoName(),
			Namespace: p.instance.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.GitRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: p.instance.Spec.ProfileURL,
			Reference: &sourcev1.GitRepositoryRef{
				Branch: p.instance.Spec.Branch,
			},
		},
	}
}

func (p *Profile) makeGitRepoName() string {
	repoParts := strings.Split(p.instance.Spec.ProfileURL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(p.instance.Name, repoName, p.instance.Spec.Branch)
}
