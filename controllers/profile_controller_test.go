package controllers_test

import (
	"context"
	"fmt"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ProfileController", func() {
	const nginxProfileURL = "https://github.com/weaveworks/nginx-profile"
	const helmChartURL = "https://charts.bitnami.com/bitnami"

	var (
		namespace            string
		branch               = "main"
		nestedArtifactBranch = "helm-artifact"
		ctx                  = context.Background()
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

	Context("Create with multiple artifacts", func() {
		DescribeTable("Applying a Profile creates the correct resources", func(pSubSpec profilesv1.ProfileSubscriptionSpec) {
			subscriptionName := "foo"

			pSub := profilesv1.ProfileSubscription{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileSubscription",
					APIVersion: "profilesubscriptions.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      subscriptionName,
					Namespace: namespace,
				},
			}
			pSub.Spec = pSubSpec
			Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

			By("creating a GitRepository resource for the profile")
			profileRepoName := "nginx-profile"
			gitRepoName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileRepoName, branch)
			gitRepo := sourcev1.GitRepository{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: gitRepoName, Namespace: namespace}, &gitRepo)
			}, 10*time.Second).ShouldNot(HaveOccurred())
			Expect(gitRepo.Spec.URL).To(Equal(nginxProfileURL))
			Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))

			By("creating a GitRepository resource for the nested profile")
			gitRepoName = fmt.Sprintf("%s-%s-%s", subscriptionName, profileRepoName, nestedArtifactBranch)
			gitRepo = sourcev1.GitRepository{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: gitRepoName, Namespace: namespace}, &gitRepo)
			}, 10*time.Second).ShouldNot(HaveOccurred())
			Expect(gitRepo.Spec.URL).To(Equal(nginxProfileURL))
			Expect(gitRepo.Spec.Reference.Branch).To(Equal(nestedArtifactBranch))

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

			By("creating a HelmRepository resource")
			helmRepoName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileRepoName, "dokuwiki")
			helmRepo := sourcev1.HelmRepository{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: helmRepoName, Namespace: namespace}, &helmRepo)
			}, 10*time.Second).ShouldNot(HaveOccurred())
			Expect(helmRepo.Spec.URL).To(Equal(helmChartURL))

			By("and creating a HelmRelease resource for the HelmRepository")
			secondProfileChartName := "dokuwiki"
			secondHelmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, secondProfileChartName)
			secondHelmRelease := helmv2.HelmRelease{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: secondHelmReleaseName, Namespace: namespace}, &secondHelmRelease)
			}, 10*time.Second).ShouldNot(HaveOccurred())
			Expect(secondHelmRelease.Spec.Chart.Spec.Chart).To(Equal("dokuwiki"))
			Expect(secondHelmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
				helmv2.CrossNamespaceObjectReference{
					Kind:      "HelmRepository",
					Name:      helmRepoName,
					Namespace: namespace,
				},
			))

			if pSub.Spec.Values != nil {
				Expect(helmRelease.Spec.Values).To(Equal(pSub.Spec.Values))
				Expect(secondHelmRelease.Spec.Values).To(Equal(pSub.Spec.Values))
			}
			if pSub.Spec.ValuesFrom != nil {
				Expect(helmRelease.Spec.ValuesFrom).To(Equal(pSub.Spec.ValuesFrom))
				Expect(secondHelmRelease.Spec.ValuesFrom).To(Equal(pSub.Spec.ValuesFrom))
			}

			conditions := []metav1.Condition{
				{
					Type:               "Ready",
					Status:             "Unknown",
					Reason:             "foo",
					Message:            "somethings wrong",
					LastTransitionTime: metav1.Now(),
				},
			}
			gitResNew := gitRepo.DeepCopyObject().(*sourcev1.GitRepository)
			gitResNew.Status.Conditions = conditions
			Expect(k8sClient.Status().Patch(ctx, gitResNew, client.MergeFrom(&gitRepo))).To(Succeed())

			helmRepNew := helmRepo.DeepCopyObject().(*sourcev1.HelmRepository)
			helmRepNew.Status.Conditions = conditions
			Expect(k8sClient.Status().Patch(ctx, helmRepNew, client.MergeFrom(&helmRepo))).To(Succeed())

			By("updating the status to Ready Unknown if the artifact resource reports it")
			profile := profilesv1.ProfileSubscription{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)).To(Succeed())
				return len(profile.Status.Conditions) > 0 &&
					profile.Status.Conditions[0].Type == "Ready" &&
					profile.Status.Conditions[0].Status == "Unknown" &&
					profile.Status.Conditions[0].Message == "somethings wrong,somethings wrong"
			}, 10*time.Second).Should(BeTrue())

			conditions = []metav1.Condition{
				{
					Type:               "Ready",
					Status:             "True",
					Reason:             "foo",
					LastTransitionTime: metav1.Now(),
				},
			}
			gitResNew = gitRepo.DeepCopyObject().(*sourcev1.GitRepository)
			gitResNew.Status.Conditions = conditions
			Expect(k8sClient.Status().Patch(ctx, gitResNew, client.MergeFrom(&gitRepo))).To(Succeed())

			helmRepNew = helmRepo.DeepCopyObject().(*sourcev1.HelmRepository)
			helmRepNew.Status.Conditions = conditions
			Expect(k8sClient.Status().Patch(ctx, helmRepNew, client.MergeFrom(&helmRepo))).To(Succeed())

			helmResNew := helmRelease.DeepCopyObject().(*helmv2.HelmRelease)
			helmResNew.Status.Conditions = conditions
			Expect(k8sClient.Status().Patch(ctx, helmResNew, client.MergeFrom(&helmRelease))).To(Succeed())

			By("updating the status to Ready True when the resources are reporting Ready")
			profile = profilesv1.ProfileSubscription{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)).To(Succeed())
				return len(profile.Status.Conditions) > 0 &&
					profile.Status.Conditions[0].Type == "Ready" &&
					profile.Status.Conditions[0].Status == "True" &&
					profile.Status.Conditions[0].Message == "all artifact resources ready"
			}, 10*time.Second).Should(BeTrue())
		},
			Entry("a single Helm chart with no supplied values", profilesv1.ProfileSubscriptionSpec{
				ProfileURL: nginxProfileURL,
				Branch:     branch,
			}),
			Entry("a single Helm chart with supplied values", profilesv1.ProfileSubscriptionSpec{
				ProfileURL: nginxProfileURL,
				Branch:     branch,
				Values: &apiextensionsv1.JSON{
					Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
				},
			}),
			Entry("a single Helm chart with values supplied via valuesFrom", profilesv1.ProfileSubscriptionSpec{
				ProfileURL: nginxProfileURL,
				Branch:     branch,
				ValuesFrom: []helmv2.ValuesReference{
					{
						Name:     "nginx-values",
						Kind:     "Secret",
						Optional: true,
					},
				},
			}),
		)

		When("retrieving the Profile Definition fails", func() {
			It("updates the status", func() {
				subscriptionName := "fetch-definition-error"
				profileURL := "https://github.com/does-not/exist"

				pSub := profilesv1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileSubscription",
						APIVersion: "profilesubscriptions.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subscriptionName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				profile := profilesv1.ProfileSubscription{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
					return err == nil && len(profile.Status.Conditions) > 0
				}, 10*time.Second, 1*time.Second).Should(BeTrue())

				Expect(profile.Status.Conditions[0].Message).To(Equal("error when fetching profile definition"))
				Expect(profile.Status.Conditions[0].Reason).To(Equal("FetchProfileFailed"))
				Expect(profile.Status.Conditions[0].Type).To(Equal("Ready"))
				Expect(profile.Status.Conditions[0].Status).To(Equal(metav1.ConditionStatus("False")))
			})
		})

		When("creating Profile artifacts fail", func() {
			It("updates the status", func() {
				subscriptionName := "git-resource-already-exists-error"
				profileURL := nginxProfileURL
				pSub := profilesv1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileSubscription",
						APIVersion: "profilesubscriptions.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subscriptionName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
						Branch:     "invalid-artifact",
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				profile := profilesv1.ProfileSubscription{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
					return err == nil && len(profile.Status.Conditions) > 0
				}, 10*time.Second, 1*time.Second).Should(BeTrue())

				Expect(profile.Status.Conditions[0].Message).To(Equal("error when reconciling profile artifacts"))
				Expect(profile.Status.Conditions[0].Type).To(Equal("Ready"))
				Expect(profile.Status.Conditions[0].Status).To(Equal(metav1.ConditionStatus("False")))
				Expect(profile.Status.Conditions[0].Reason).To(Equal("CreateFailed"))
			})
		})
	})
})
