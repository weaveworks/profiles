package gitrepository_test

import (
	"context"
	"fmt"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/gitrepository"
	"github.com/weaveworks/profiles/pkg/gitrepository/fakes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Gitrepository", func() {
	var (
		manager   *gitrepository.Manager
		kClient   *fakes.FakeKubernetes
		callCount int
		repo      = profilesv1.Repository{
			URL: "github.com/example/repo",
			SecretRef: &meta.LocalObjectReference{
				Name: "my-secret",
			},
		}
	)

	BeforeEach(func() {
		callCount = 0
		kClient = new(fakes.FakeKubernetes)
		manager = gitrepository.NewManager(context.TODO(), "profiles-system", kClient, time.Second, time.Millisecond)
	})

	Describe("CreateAndWaitForResources", func() {

		When("the gitrepositorys create successfully", func() {
			BeforeEach(func() {
				kClient.GetStub = func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
					callCount++
					if callCount == 3 {
						Expect(obj.(*sourcev1.GitRepository).Name).To(Equal("repo-v0.1.0"))
						obj.(*sourcev1.GitRepository).Status = sourcev1.GitRepositoryStatus{
							URL: "url1",
						}
					}
					if callCount == 4 {
						Expect(obj.(*sourcev1.GitRepository).Name).To(Equal("repo-foo-v1.0.0"))
						obj.(*sourcev1.GitRepository).Status = sourcev1.GitRepositoryStatus{
							URL: "url2",
						}
					}
					return nil
				}
			})

			It("creates the gitrepository resources and waits for them to be ready", func() {
				resources, err := manager.CreateAndWaitForResources(repo, []gitrepository.Instance{
					{
						Tag:  "v0.1.0",
						Path: "profile.yaml",
					},
					{
						Tag:  "foo/v1.0.0",
						Path: "foo/profile.yaml",
					},
				})

				Expect(err).NotTo(HaveOccurred())
				ignore1 := `# exclude all
/*
# include deploy dir
!/profile.yaml`

				ignore2 := `# exclude all
/*
# include deploy dir
!/foo/profile.yaml`
				Expect(kClient.CreateCallCount()).To(Equal(2))
				Expect(kClient.GetCallCount()).To(Equal(4))

				Expect(resources).To(ConsistOf(
					&sourcev1.GitRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "repo-v0.1.0",
							Namespace: "profiles-system",
						},
						TypeMeta: metav1.TypeMeta{
							Kind:       sourcev1.GitRepositoryKind,
							APIVersion: sourcev1.GroupVersion.String(),
						},
						Spec: sourcev1.GitRepositorySpec{
							URL: "github.com/example/repo",
							Reference: &sourcev1.GitRepositoryRef{
								Tag: "v0.1.0",
							},
							Ignore: &ignore1,
							SecretRef: &meta.LocalObjectReference{
								Name: "my-secret",
							},
						},
						Status: sourcev1.GitRepositoryStatus{
							URL: "url1",
						},
					},
					&sourcev1.GitRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "repo-foo-v1.0.0",
							Namespace: "profiles-system",
						},
						TypeMeta: metav1.TypeMeta{
							Kind:       sourcev1.GitRepositoryKind,
							APIVersion: sourcev1.GroupVersion.String(),
						},
						Spec: sourcev1.GitRepositorySpec{
							URL: "github.com/example/repo",
							Reference: &sourcev1.GitRepositoryRef{
								Tag: "foo/v1.0.0",
							},
							Ignore: &ignore2,
							SecretRef: &meta.LocalObjectReference{
								Name: "my-secret",
							},
						},
						Status: sourcev1.GitRepositoryStatus{
							URL: "url2",
						},
					},
				))
			})
		})

		When("create fails", func() {
			BeforeEach(func() {
				kClient.CreateReturns(fmt.Errorf("createfailed"))
			})

			It("returns an error", func() {
				_, err := manager.CreateAndWaitForResources(repo, []gitrepository.Instance{
					{
						Tag:  "v0.1.0",
						Path: "",
					},
					{
						Tag:  "foo/v1.0.0",
						Path: "foo",
					},
				})

				Expect(err).To(MatchError("failed to create gitrepository: createfailed"))
			})
		})

		When("get fails", func() {
			BeforeEach(func() {
				kClient.GetReturns(fmt.Errorf("getfailed"))
			})

			It("returns an error", func() {
				_, err := manager.CreateAndWaitForResources(repo, []gitrepository.Instance{
					{
						Tag:  "v0.1.0",
						Path: "",
					},
					{
						Tag:  "foo/v1.0.0",
						Path: "foo",
					},
				})

				Expect(err).To(MatchError("failed to get gitrepository: getfailed"))
			})
		})

		When("timesout waiting for status to change", func() {
			It("returns an error", func() {
				_, err := manager.CreateAndWaitForResources(repo, []gitrepository.Instance{
					{
						Tag:  "v0.1.0",
						Path: "",
					},
					{
						Tag:  "foo/v1.0.0",
						Path: "foo",
					},
				})

				Expect(err).To(MatchError("timed out waiting for profiles-system/repo-v0.1.0 gitrepository.Status.URL to be populated"))
			})
		})
	})

	Describe("DeleteResources", func() {
		It("deletes all resoures", func() {
			resources := []*sourcev1.GitRepository{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-v0.1.0",
						Namespace: "profiles-system",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       sourcev1.GitRepositoryKind,
						APIVersion: sourcev1.GroupVersion.String(),
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "v0.1.0",
						},
						SecretRef: &meta.LocalObjectReference{
							Name: "my-secret",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "url1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "repo-foo-v1.0.0",
						Namespace: "profiles-system",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       sourcev1.GitRepositoryKind,
						APIVersion: sourcev1.GroupVersion.String(),
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: "github.com/example/repo",
						Reference: &sourcev1.GitRepositoryRef{
							Tag: "foo/v1.0.0",
						},
						SecretRef: &meta.LocalObjectReference{
							Name: "my-secret",
						},
					},
					Status: sourcev1.GitRepositoryStatus{
						URL: "url2",
					},
				},
			}

			Expect(manager.DeleteResources(resources)).To(Succeed())
			Expect(kClient.DeleteCallCount()).To(Equal(2))
			_, res, _ := kClient.DeleteArgsForCall(0)
			Expect(res.(*sourcev1.GitRepository)).To(Equal(resources[0]))
			_, res, _ = kClient.DeleteArgsForCall(1)
			Expect(res.(*sourcev1.GitRepository)).To(Equal(resources[1]))
		})
		When("delete fails", func() {
			It("returns an error", func() {
				resources := []*sourcev1.GitRepository{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "repo-v0.1.0",
							Namespace: "profiles-system",
						},
						TypeMeta: metav1.TypeMeta{
							Kind:       sourcev1.GitRepositoryKind,
							APIVersion: sourcev1.GroupVersion.String(),
						},
						Spec: sourcev1.GitRepositorySpec{
							URL: "github.com/example/repo",
							Reference: &sourcev1.GitRepositoryRef{
								Tag: "v0.1.0",
							},
							SecretRef: &meta.LocalObjectReference{
								Name: "my-secret",
							},
						},
						Status: sourcev1.GitRepositoryStatus{
							URL: "url1",
						},
					},
				}

				kClient.DeleteReturns(fmt.Errorf("foo"))
				Expect(manager.DeleteResources(resources)).To(MatchError("failed to delete resource profiles-system/repo-v0.1.0: foo"))
			})
		})
	})
})
