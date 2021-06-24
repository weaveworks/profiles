package git

import (
	"fmt"

	"github.com/fluxcd/source-controller/pkg/git/gogit"
	extgogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	corev1 "k8s.io/api/core/v1"
)

//Client git client
type Client struct{}

//ListTags returns a list of tags for a given repository
func (c *Client) ListTags(url string, secret *corev1.Secret) ([]string, error) {
	rem := extgogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	var auth transport.AuthMethod
	if secret != nil {
		authStrategy, err := gogit.AuthSecretStrategyForURL(url)

		if err != nil {
			return nil, fmt.Errorf("failed to create auth strateg from URL %q : %w", url, err)
		}
		authMethod, err := authStrategy.Method(*secret)
		if err != nil {
			return nil, fmt.Errorf("failed to get auth method: %w", err)
		}
		auth = authMethod.AuthMethod
	}

	refs, err := rem.List(&extgogit.ListOptions{
		Auth: auth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	return tags, nil
}
