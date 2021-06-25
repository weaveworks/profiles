package scanner_test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/gitrepository"
	"github.com/weaveworks/profiles/pkg/scanner"
	"github.com/weaveworks/profiles/pkg/scanner/fakes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Scanner", func() {
	var (
		s              scanner.RepoScanner
		gitClient      *fakes.FakeGitClient
		gitRepoManager *fakes.FakeGitRepositoryManager
		httpClient     *fakes.FakeHTTPClient
		repoSecret     = &corev1.Secret{
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
		}
		repo = profilesv1.Repository{
			URL: "github.com/example/repo",
			SecretRef: &meta.LocalObjectReference{
				Name: "foo",
			},
		}
	)

	BeforeEach(func() {
		gitClient = new(fakes.FakeGitClient)
		gitRepoManager = new(fakes.FakeGitRepositoryManager)
		httpClient = new(fakes.FakeHTTPClient)
		s = scanner.New(gitRepoManager, gitClient, httpClient, logr.Discard())
	})

	Context("when the repo has matching tags", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns([]string{"name/v0.0.1", "name/v0.1.0", "v1.0.0", "some-notsemver"}, nil)
			gitRepoManager.CreateAndWaitForResourcesReturns([]*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.one",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-foo-v1.0.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "foo/v1.0.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.two",
					},
				},
			}, nil)

			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusOK,
				Body: tarContents([]byte(`---
spec:
  name: other-name
  description: some desc
  maintainer: me
  Prerequisites:
  - stuff`))}, nil)
			httpClient.DoReturnsOnCall(1, &http.Response{
				StatusCode: http.StatusOK,
				Body: tarContents([]byte(`---
spec:
  name: foo-name
  description: some desc
  maintainer: me
  Prerequisites:
  - stuff`))}, nil)
		})

		It("returns a list of profiles", func() {
			profiles, tags, err := s.ScanRepository(repo, repoSecret, []string{"name/v0.0.1"})
			Expect(err).NotTo(HaveOccurred())

			Expect(gitClient.ListTagsCallCount()).To(Equal(1))
			url, secret := gitClient.ListTagsArgsForCall(0)
			Expect(url).To(Equal("github.com/example/repo"))
			Expect(secret).To(Equal(repoSecret))

			Expect(gitRepoManager.CreateAndWaitForResourcesCallCount()).To(Equal(1))
			givenRepo, repos := gitRepoManager.CreateAndWaitForResourcesArgsForCall(0)
			Expect(givenRepo).To(Equal(repo))
			Expect(url).To(Equal("github.com/example/repo"))
			Expect(secret).To(Equal(repoSecret))
			Expect(repos).To(ConsistOf(
				gitrepository.Instance{
					Tag:  "name/v0.1.0",
					Path: "name/profile.yaml",
				},
				gitrepository.Instance{
					Tag:  "v1.0.0",
					Path: "profile.yaml",
				},
			))
			Expect(httpClient.DoCallCount()).To(Equal(2))
			Expect(httpClient.DoArgsForCall(0).URL.String()).To(Equal("tarball.one"))
			Expect(httpClient.DoArgsForCall(1).URL.String()).To(Equal("tarball.two"))

			Expect(gitRepoManager.DeleteResourcesCallCount()).To(Equal(1))
			Expect(gitRepoManager.DeleteResourcesArgsForCall(0)).To(ConsistOf(
				&sourcev1.GitRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.one",
					},
				},
				&sourcev1.GitRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-foo-v1.0.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "foo/v1.0.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.two",
					},
				},
			))

			Expect(profiles).To(ConsistOf(profilesv1.ProfileCatalogEntry{
				ProfileDescription: profilesv1.ProfileDescription{
					Name:          "foo-name",
					Description:   "some desc",
					Maintainer:    "me",
					Prerequisites: []string{"stuff"},
				},
				Tag: "foo/v1.0.0",
				URL: "github.com/example/repo",
			}, profilesv1.ProfileCatalogEntry{
				ProfileDescription: profilesv1.ProfileDescription{
					Name:          "other-name",
					Description:   "some desc",
					Maintainer:    "me",
					Prerequisites: []string{"stuff"},
				},
				Tag: "v0.1.0",
				URL: "github.com/example/repo",
			}))
			Expect(tags).To(ConsistOf("name/v0.1.0", "v1.0.0", "some-notsemver"))
		})
	})

	When("ListTags fails", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns(nil, fmt.Errorf("listfail"))
		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError("failed to list tags: listfail"))

		})
	})

	When("CreateAndWaitForResources fails", func() {
		BeforeEach(func() {
			gitRepoManager.CreateAndWaitForResourcesReturns(nil, fmt.Errorf("createfail"))
		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError("failed to create gitrepository resources: createfail"))
		})
	})

	When("the tarball url is invalid", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns([]string{"name/v0.1.0", "v1.0.0", "some-notsemver"}, nil)
			gitRepoManager.CreateAndWaitForResourcesReturns([]*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "invaludurl{DEf1=ghi@example.com:5432/db?sslmode=require",
					},
				},
			}, nil)

		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError(ContainSubstring("failed to create request:")))
		})
	})

	When("httpclient.Do fails", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns([]string{"name/v0.1.0", "v1.0.0", "some-notsemver"}, nil)
			gitRepoManager.CreateAndWaitForResourcesReturns([]*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.one",
					},
				},
			}, nil)

			httpClient.DoReturnsOnCall(0, &http.Response{}, fmt.Errorf("dofail"))
		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError("failed to GET \"tarball.one\": dofail"))
		})
	})

	When("request returns non 200", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns([]string{"name/v0.1.0", "v1.0.0", "some-notsemver"}, nil)
			gitRepoManager.CreateAndWaitForResourcesReturns([]*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.one",
					},
				},
			}, nil)

			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       gbytes.NewBuffer(),
			}, nil)
		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError("request failed status code 400"))
		})
	})

	When("the body isn't a tarball", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns([]string{"name/v0.1.0", "v1.0.0", "some-notsemver"}, nil)
			gitRepoManager.CreateAndWaitForResourcesReturns([]*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.one",
					},
				},
			}, nil)

			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusOK,
				Body:       gbytes.NewBuffer(),
			}, nil)
		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError(ContainSubstring("failed to parse tarball:")))
		})
	})

	When("the file isn't valid yaml", func() {
		BeforeEach(func() {
			gitClient.ListTagsReturns([]string{"name/v0.1.0"}, nil)
			gitRepoManager.CreateAndWaitForResourcesReturns([]*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "tarball.one",
					},
				},
			}, nil)

			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusOK,
				Body:       tarContents([]byte(`!@\:1\23notyaml`))}, nil)
		})

		It("returns an error", func() {
			_, _, err := s.ScanRepository(repo, repoSecret, nil)
			Expect(err).To(MatchError(ContainSubstring("failed to decode profile.yaml:")))
		})
	})
})

func tarContents(content []byte) io.ReadCloser {
	buf := gbytes.NewBuffer()
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)
	hdr := &tar.Header{
		Name: "profile.yaml",
		Mode: 0600,
		Size: int64(len(content)),
	}
	Expect(tw.WriteHeader(hdr)).To(Succeed())
	_, err := tw.Write([]byte(content))
	Expect(err).NotTo(HaveOccurred())
	Expect(gw.Close()).To(Succeed())
	return buf
}
