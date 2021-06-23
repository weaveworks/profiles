package scanner

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fluxcd/pkg/version"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/gitrepository"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_git_client.go . GitClient
//GitClient client for interacting with git
type GitClient interface {
	ListTags(url string, secret *corev1.Secret) ([]string, error)
}

//counterfeiter:generate -o fakes/fake_repo_manager.go . GitRepositoryManager
//GitRepositoryManager for managing gitrepositorys
type GitRepositoryManager interface {
	CreateAndWaitForResources(repo profilesv1.Repository, tags []gitrepository.Instance) ([]*sourcev1.GitRepository, error)
	DeleteResources([]*sourcev1.GitRepository) error
}

//counterfeiter:generate -o fakes/fake_http_client.go . HTTPClient
//HTTPClient for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

//Scanner for scanning repositorys
type Scanner struct {
	gitRepositoryManager GitRepositoryManager
	gitClient            GitClient
	httpClient           HTTPClient
	logger               logr.Logger
}

//New returns a Scanner
func New(gitRepositoryManager GitRepositoryManager, gitClient GitClient, httpClient HTTPClient, logger logr.Logger) *Scanner {
	return &Scanner{
		gitRepositoryManager: gitRepositoryManager,
		gitClient:            gitClient,
		httpClient:           httpClient,
		logger:               logger,
	}
}

//ScanRepository for profiles
func (s *Scanner) ScanRepository(repo profilesv1.Repository, secret *corev1.Secret) ([]profilesv1.ProfileCatalogEntry, error) {
	tags, err := s.gitClient.ListTags(repo.URL, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	s.logger.Info("found tags", "url", repo.URL, "tags", tags)

	var instances []gitrepository.Instance
	for _, tag := range tags {
		semver, path := getSemverAndPathFromTag(tag)
		if _, err := version.ParseVersion(semver); err == nil {
			instances = append(instances, gitrepository.Instance{
				Tag:  tag,
				Path: path,
			})
		}
	}

	gitRepositoryResources, err := s.gitRepositoryManager.CreateAndWaitForResources(repo, instances)
	if err != nil {
		return nil, fmt.Errorf("failed to create gitrepository resources: %w", err)
	}
	s.logger.Info("gitrepositorys created", "gitrepositories", gitRepositoryResources)

	defer func() {
		if err := s.gitRepositoryManager.DeleteResources(gitRepositoryResources); err != nil {
			s.logger.Error(err, "failed to cleanup git resources", "gitrepositories", gitRepositoryResources)
		}
	}()

	var profiles []profilesv1.ProfileCatalogEntry
	for _, gitRepo := range gitRepositoryResources {
		profileDef, err := s.fetchProfileFromTarball(gitRepo)
		if err != nil {
			return nil, err
		}
		if profileDef != nil && profileDef.Spec.Name != "" {
			profiles = append(profiles, profilesv1.ProfileCatalogEntry{
				ProfileDescription: profileDef.Spec.ProfileDescription,
				Tag:                gitRepo.Spec.Reference.Tag,
				URL:                repo.URL,
			})
		}
	}

	return profiles, nil
}

func (s *Scanner) fetchProfileFromTarball(gitRepo *sourcev1.GitRepository) (*profilesv1.ProfileDefinition, error) {
	req, err := http.NewRequest("GET", gitRepo.Status.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to GET %q: %w", gitRepo.Status.URL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed status code %d", resp.StatusCode)
	}

	return extractProfileFromTarball(resp.Body)
}

func extractProfileFromTarball(gzipStream io.Reader) (*profilesv1.ProfileDefinition, error) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tarball: %w", err)
	}
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			return nil, nil
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read tarball file: %w", err)
		}

		if header.Typeflag == tar.TypeReg {
			decoder := yaml.NewYAMLOrJSONDecoder(tarReader, 10000)
			var profileDef profilesv1.ProfileDefinition
			if err = decoder.Decode(&profileDef); err != nil {
				return nil, fmt.Errorf("failed to decode profile.yaml: %w", err)
			}
			return &profileDef, nil
		}
	}
}

func getSemverAndPathFromTag(tag string) (string, string) {
	v := tag
	path := "profile.yaml"
	splitTag := strings.Split(tag, "/")
	if len(splitTag) == 2 {
		path = filepath.Join(splitTag[0], path)
		v = splitTag[1]
	}
	return v, path
}
