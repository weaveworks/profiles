package controllers_test

import (
	"context"
	"fmt"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/scanner"
	"github.com/weaveworks/profiles/pkg/scanner/fakes"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ProfileCatalogSourceController", func() {
	var (
		namespace string
		ctx       = context.Background()
	)

	BeforeEach(func() {
		namespace = uuid.New().String()
		nsp := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		Expect(k8sClient.Create(context.Background(), &nsp)).To(Succeed())
	})

	When("providing a static list of profiles", func() {
		It("syncs the in-memory list when a ProfileCatalogSource is added or deleted", func() {
			By("creating a new ProfileCatalogSource")
			catalogSource := &profilesv1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: "profile.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "catalog",
					Namespace: namespace,
				},
				Spec: profilesv1.ProfileCatalogSourceSpec{
					Profiles: []profilesv1.ProfileCatalogEntry{
						{
							ProfileDescription: profilesv1.ProfileDescription{
								Description: "bar",
							},
							Name: "foo",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, catalogSource)).Should(Succeed())

			By("searching for a profile")
			query := func() []profilesv1.ProfileCatalogEntry {
				return catalogReconciler.Profiles.Search("foo")
			}
			Eventually(query, 2*time.Second).Should(ContainElement(profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Description: "bar"}, Name: "foo", CatalogSource: "catalog"}))
			Expect(k8sClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: "catalog"}, catalogSource)).To(Succeed())

			By("adding more items to ProfileCatalogSource")
			pName := fmt.Sprintf("new-profile-%s", uuid.New().String())
			catalogSource.Spec.Profiles = append(catalogSource.Spec.Profiles, profilesv1.ProfileCatalogEntry{
				ProfileDescription: profilesv1.ProfileDescription{
					Description: "I am new here",
				},
				Name: pName,
			})
			Expect(k8sClient.Update(context.Background(), catalogSource)).To(Succeed())

			Eventually(func() []profilesv1.ProfileCatalogEntry {
				return catalogReconciler.Profiles.Search(pName)
			}, 2*time.Second).Should(ConsistOf(profilesv1.ProfileCatalogEntry{
				ProfileDescription: profilesv1.ProfileDescription{
					Description: "I am new here",
				},
				Name:          pName,
				CatalogSource: "catalog",
			}))

			By("deleting the ProfileCatalogSource")
			Expect(k8sClient.Delete(ctx, catalogSource)).To(Succeed())
			Eventually(query, 2*time.Second).Should(BeEmpty())
			Expect(catalogReconciler.Profiles.Search(pName)).To(BeEmpty())
		})
	})

	When("providing a repo to scan", func() {
		var catalogSource *profilesv1.ProfileCatalogSource
		BeforeEach(func() {
			fakeRepoScanner = new(fakes.FakeRepoScanner)
			catalogReconciler.SetNewScanner(
				func(gitRepositoryManager scanner.GitRepositoryManager, gitClient scanner.GitClient, httpClients scanner.HTTPClient, logger logr.Logger) scanner.RepoScanner {
					return fakeRepoScanner
				},
			)
			fakeRepoScanner.ScanRepositoryReturnsOnCall(0, []profilesv1.ProfileCatalogEntry{
				{
					Name: "foo",
				},
			}, []string{"foo"}, nil)

			By("creating a new ProfileCatalogSource")
			catalogSource = &profilesv1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: "profile.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "catalog-2",
					Namespace: namespace,
				},
				Spec: profilesv1.ProfileCatalogSourceSpec{
					Repos: []profilesv1.Repository{
						{
							URL: "github.com/weaveworks/profiles-examples",
							SecretRef: &meta.LocalObjectReference{
								Name: "my-secret",
							},
						},
					},
				},
			}
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-secret",
					Namespace: namespace,
				},
				Data: map[string][]byte{},
			}
			Expect(k8sClient.Create(ctx, secret)).Should(Succeed())
			Expect(k8sClient.Create(ctx, catalogSource)).Should(Succeed())
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, catalogSource)).Should(Succeed())
			catalogReconciler.Profiles.Remove("catalog-2")
		})

		It("scans the repository", func() {
			By("searching for a profile")
			query := func() []profilesv1.ProfileCatalogEntry {
				return catalogReconciler.Profiles.Search("foo")
			}
			Eventually(query, 2*time.Second).Should(ContainElement(profilesv1.ProfileCatalogEntry{Name: "foo", CatalogSource: "catalog-2"}))

			By("only searching for new tags")
			Eventually(func() int {
				return fakeRepoScanner.ScanRepositoryCallCount()
			}).Should(Equal(2))
			repo, secret, tags := fakeRepoScanner.ScanRepositoryArgsForCall(0)
			Expect(repo).To(Equal(profilesv1.Repository{URL: "github.com/weaveworks/profiles-examples", SecretRef: &meta.LocalObjectReference{Name: "my-secret"}}))
			Expect(secret.Name).To(Equal("my-secret"))
			Expect(tags).To(BeNil())

			repo, secret, tags = fakeRepoScanner.ScanRepositoryArgsForCall(1)
			Expect(repo).To(Equal(profilesv1.Repository{URL: "github.com/weaveworks/profiles-examples", SecretRef: &meta.LocalObjectReference{Name: "my-secret"}}))
			Expect(secret.Name).To(Equal("my-secret"))
			Expect(tags).To(ConsistOf("foo"))

			Expect(k8sClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: "catalog-2"}, catalogSource)).To(Succeed())
			Expect(catalogSource.Status.ScannedRepositories).To(ConsistOf(
				profilesv1.ScannedRepository{
					URL:  "github.com/weaveworks/profiles-examples",
					Tags: []string{"foo"},
				},
			))
		})

		When("the catalog gets wiped", func() {
			It("re-scans the repository, resetting the tags on the status", func() {
				By("searching for a profile")
				query := func() []profilesv1.ProfileCatalogEntry {
					return catalogReconciler.Profiles.Search("foo")
				}
				Eventually(query, 2*time.Second).Should(ContainElement(profilesv1.ProfileCatalogEntry{Name: "foo", CatalogSource: "catalog-2"}))

				Eventually(func() []profilesv1.ScannedRepository {
					Expect(k8sClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: "catalog-2"}, catalogSource)).To(Succeed())
					return catalogSource.Status.ScannedRepositories
				}, 2*time.Second).Should(ConsistOf(
					profilesv1.ScannedRepository{
						URL:  "github.com/weaveworks/profiles-examples",
						Tags: []string{"foo"},
					},
				))
				Expect(fakeRepoScanner.ScanRepositoryCallCount()).To(Equal(2))

				By("rescanning the repository when the catalog gets reset")
				fakeRepoScanner.ScanRepositoryReturnsOnCall(2, []profilesv1.ProfileCatalogEntry{
					{
						Name: "bar",
					},
					{
						Name: "baz",
					},
				}, []string{"bar", "baz"}, nil)

				catalogReconciler.Profiles.Remove("catalog-2")
				//force a reconciliation loop
				catalogSource.Labels = map[string]string{"some": "label"}
				Expect(k8sClient.Update(ctx, catalogSource)).Should(Succeed())

				Eventually(func() int {
					return fakeRepoScanner.ScanRepositoryCallCount()
				}, time.Second*2).Should(Equal(4))
				repo, secret, tags := fakeRepoScanner.ScanRepositoryArgsForCall(2)
				Expect(repo).To(Equal(profilesv1.Repository{URL: "github.com/weaveworks/profiles-examples", SecretRef: &meta.LocalObjectReference{Name: "my-secret"}}))
				Expect(secret.Name).To(Equal("my-secret"))
				Expect(tags).To(BeNil())

				query = func() []profilesv1.ProfileCatalogEntry {
					return catalogReconciler.Profiles.Search("baz")
				}
				Eventually(query, 2*time.Second).Should(ContainElement(profilesv1.ProfileCatalogEntry{Name: "baz", CatalogSource: "catalog-2"}))
				Eventually(func() []profilesv1.ScannedRepository {
					Expect(k8sClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: "catalog-2"}, catalogSource)).To(Succeed())
					return catalogSource.Status.ScannedRepositories
				}, 2*time.Second).Should(ConsistOf(
					profilesv1.ScannedRepository{
						URL:  "github.com/weaveworks/profiles-examples",
						Tags: []string{"bar", "baz"},
					},
				))

				repo, secret, tags = fakeRepoScanner.ScanRepositoryArgsForCall(2)
				Expect(repo).To(Equal(profilesv1.Repository{URL: "github.com/weaveworks/profiles-examples", SecretRef: &meta.LocalObjectReference{Name: "my-secret"}}))
				Expect(secret.Name).To(Equal("my-secret"))
				Expect(tags).To(BeNil())

				repo, secret, tags = fakeRepoScanner.ScanRepositoryArgsForCall(3)
				Expect(repo).To(Equal(profilesv1.Repository{URL: "github.com/weaveworks/profiles-examples", SecretRef: &meta.LocalObjectReference{Name: "my-secret"}}))
				Expect(secret.Name).To(Equal("my-secret"))
				Expect(tags).To(ConsistOf("bar", "baz"))
			})
		})
	})
})
