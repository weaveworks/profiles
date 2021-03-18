package git

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-logr/logr"
	"github.com/weaveworks/profiles/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:generate counterfeiter -o fakes/fake_http_client.go . HTTPClient
type HTTPClient interface {
	Get(string) (*http.Response, error)
}

var httpClient HTTPClient = http.DefaultClient

func GetProfileDefinition(repoURL, branch string, log logr.Logger) (v1alpha1.ProfileDefinition, error) {
	if _, err := url.Parse(repoURL); err != nil {
		return v1alpha1.ProfileDefinition{}, fmt.Errorf("failed to parse repo URL %q: %w", repoURL, err)
	}

	if !strings.Contains(repoURL, "github.com") {
		return v1alpha1.ProfileDefinition{}, errors.New("unsupported git provider, only github.com is currently supported")
	}

	profileURL := strings.Replace(repoURL, "github.com", "raw.githubusercontent.com", 1)
	profileURL = fmt.Sprintf("%s/%s/profile.yaml", profileURL, branch)

	log.Info("fetching profile.yaml", "repoURL", repoURL)
	resp, err := httpClient.Get(profileURL)
	if err != nil {
		return v1alpha1.ProfileDefinition{}, fmt.Errorf("failed to fetch profile: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return v1alpha1.ProfileDefinition{}, fmt.Errorf("failed to fetch profile: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	profile := v1alpha1.ProfileDefinition{}
	err = yaml.NewYAMLOrJSONDecoder(resp.Body, 4096).Decode(&profile)
	if err != nil {
		return v1alpha1.ProfileDefinition{}, fmt.Errorf("failed to parse profile: %w", err)
	}

	return profile, nil
}
