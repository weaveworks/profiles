package controllers_test

import (
	"context"
	"fmt"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ProfileController", func() {
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

	Context("Create", func() {
		When("the Profile contains a single local Helm Chart", func() {
			It("creates the correct resources", func() {
				subscriptionName := "foo"
				branch := "main"
				profileURL := "https://github.com/weaveworks/nginx-profile"

				pSub := v1alpha1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileSubscription",
						APIVersion: "profilesubscriptions.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subscriptionName,
						Namespace: namespace,
					},
					Spec: v1alpha1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				By("creating a GitRepository resource")
				profileRepoName := "nginx-profile"
				gitRepoName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileRepoName, branch)
				gitRepo := sourcev1.GitRepository{}
				Eventually(func() error {
					return k8sClient.Get(ctx, client.ObjectKey{Name: gitRepoName, Namespace: namespace}, &gitRepo)
				}, 10*time.Second).ShouldNot(HaveOccurred())
				Expect(gitRepo.Spec.URL).To(Equal(profileURL))
				Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))

				By("creating a HelmRelease resource")
				profileName := "nginx"
				chartName := "nginx-server"
				helmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName)
				helmRelease := helmv2.HelmRelease{}
				Eventually(func() error {
					return k8sClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
				}, 10*time.Second).ShouldNot(HaveOccurred())
				Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal("nginx/chart"))
				Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
					helmv2.CrossNamespaceObjectReference{
						Kind:      "GitRepository",
						Name:      gitRepoName,
						Namespace: namespace,
					},
				))

				By("updating the status")
				profile := v1alpha1.ProfileSubscription{}
				Eventually(func() string {
					Expect(k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)).To(Succeed())
					return profile.Status.State
				}, 10*time.Second).Should(Equal("running"))
			})
		})

		When("the retrieving the Profile Definition fails", func() {
			It("updates the status", func() {
				subscriptionName := "fetch-definition-error"
				profileURL := "https://github.com/does-not/exist"

				pSub := v1alpha1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileSubscription",
						APIVersion: "profilesubscriptions.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subscriptionName,
						Namespace: namespace,
					},
					Spec: v1alpha1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				profile := v1alpha1.ProfileSubscription{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
					return err == nil && profile.Status != v1alpha1.ProfileSubscriptionStatus{}
				}, 10*time.Second, 1*time.Second).Should(BeTrue())

				Expect(profile.Status.Message).To(Equal("error when fetching profile definition"))
				Expect(profile.Status.State).To(Equal("failing"))
			})
		})

		When("the creating Profile artifacts fail", func() {
			It("updates the status", func() {
				subscriptionName := "git-resource-already-exists-error"
				profileURL := "https://github.com/weaveworks/nginx-profile"

				gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "nginx-profile", "main")
				gitRepo := sourcev1.GitRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      gitRefName,
						Namespace: namespace,
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "GitRepository",
						APIVersion: "source.toolkit.fluxcd.io/v1beta1",
					},
					Spec: sourcev1.GitRepositorySpec{
						URL: profileURL,
						Reference: &sourcev1.GitRepositoryRef{
							Branch: "main",
						},
					},
				}
				Expect(k8sClient.Create(ctx, &gitRepo)).Should(Succeed())

				pSub := v1alpha1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileSubscription",
						APIVersion: "profilesubscriptions.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subscriptionName,
						Namespace: namespace,
					},
					Spec: v1alpha1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				profile := v1alpha1.ProfileSubscription{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
					return err == nil && profile.Status != v1alpha1.ProfileSubscriptionStatus{}
				}, 10*time.Second, 1*time.Second).Should(BeTrue())

				Expect(profile.Status.Message).To(Equal("error when creating profile artifacts"))
				Expect(profile.Status.State).To(Equal("failing"))
			})
		})
	})
})
