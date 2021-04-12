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

	//Context("Create Local", func() {
	//	DescribeTable("Applying a Profile creates the correct resources", func(pSubSpec profilesv1.ProfileSubscriptionSpec) {
	//		subscriptionName := "foo"
	//		branch := pSubSpec.Branch
	//
	//		pSub := profilesv1.ProfileSubscription{
	//			TypeMeta: metav1.TypeMeta{
	//				Kind:       "ProfileSubscription",
	//				APIVersion: "profilesubscriptions.weave.works/v1alpha1",
	//			},
	//			ObjectMeta: metav1.ObjectMeta{
	//				Name:      subscriptionName,
	//				Namespace: namespace,
	//			},
	//		}
	//		pSub.Spec = pSubSpec
	//		Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())
	//
	//		By("creating a GitRepository resource")
	//		profileRepoName := "nginx-profile"
	//		gitRepoName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileRepoName, branch)
	//		gitRepo := sourcev1.GitRepository{}
	//		Eventually(func() error {
	//			return k8sClient.Get(ctx, client.ObjectKey{Name: gitRepoName, Namespace: namespace}, &gitRepo)
	//		}, 10*time.Second).ShouldNot(HaveOccurred())
	//		Expect(gitRepo.Spec.URL).To(Equal(nginxProfileURL))
	//		Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))
	//
	//		By("creating a HelmRelease resource")
	//		profileName := "nginx"
	//		chartName := "nginx-server"
	//		helmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName)
	//		helmRelease := helmv2.HelmRelease{}
	//		Eventually(func() error {
	//			return k8sClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
	//		}, 10*time.Second).ShouldNot(HaveOccurred())
	//		Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal("nginx/chart"))
	//		Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
	//			helmv2.CrossNamespaceObjectReference{
	//				Kind:      "GitRepository",
	//				Name:      gitRepoName,
	//				Namespace: namespace,
	//			},
	//		))
	//		if pSub.Spec.Values != nil {
	//			Expect(helmRelease.Spec.Values).To(Equal(pSub.Spec.Values))
	//		}
	//		if pSub.Spec.ValuesFrom != nil {
	//			Expect(helmRelease.Spec.ValuesFrom).To(Equal(pSub.Spec.ValuesFrom))
	//		}
	//
	//		By("updating the status")
	//		profile := profilesv1.ProfileSubscription{}
	//		Eventually(func() string {
	//			Expect(k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)).To(Succeed())
	//			return profile.Status.State
	//		}, 10*time.Second).Should(Equal("running"))
	//	},
	//		Entry("a single Helm chart with no supplied values", profilesv1.ProfileSubscriptionSpec{
	//			ProfileURL: nginxProfileURL,
	//			Branch:     "support-helm-urls",
	//		}),
	//		Entry("a single Helm chart with supplied values", profilesv1.ProfileSubscriptionSpec{
	//			ProfileURL: nginxProfileURL,
	//			Branch:     "support-helm-urls",
	//			Values: &apiextensionsv1.JSON{
	//				Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
	//			},
	//		}),
	//		Entry("a single Helm chart with values supplied via valuesFrom", profilesv1.ProfileSubscriptionSpec{
	//			ProfileURL: nginxProfileURL,
	//			Branch:     "support-helm-urls",
	//			ValuesFrom: []helmv2.ValuesReference{
	//				{
	//					Name:     "nginx-values",
	//					Kind:     "Secret",
	//					Optional: true,
	//				},
	//			},
	//		}),
	//	)
	//
	//	When("retrieving the Profile Definition fails", func() {
	//		It("updates the status", func() {
	//			subscriptionName := "fetch-definition-error"
	//			profileURL := "https://github.com/does-not/exist"
	//
	//			pSub := profilesv1.ProfileSubscription{
	//				TypeMeta: metav1.TypeMeta{
	//					Kind:       "ProfileSubscription",
	//					APIVersion: "profilesubscriptions.weave.works/v1alpha1",
	//				},
	//				ObjectMeta: metav1.ObjectMeta{
	//					Name:      subscriptionName,
	//					Namespace: namespace,
	//				},
	//				Spec: profilesv1.ProfileSubscriptionSpec{
	//					ProfileURL: profileURL,
	//				},
	//			}
	//			Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())
	//
	//			profile := profilesv1.ProfileSubscription{}
	//			Eventually(func() bool {
	//				err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
	//				return err == nil && profile.Status != profilesv1.ProfileSubscriptionStatus{}
	//			}, 10*time.Second, 1*time.Second).Should(BeTrue())
	//
	//			Expect(profile.Status.Message).To(Equal("error when fetching profile definition"))
	//			Expect(profile.Status.State).To(Equal("failing"))
	//		})
	//	})
	//
	//	When("creating Profile artifacts fail", func() {
	//		It("updates the status", func() {
	//			subscriptionName := "git-resource-already-exists-error"
	//			profileURL := nginxProfileURL
	//
	//			gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "nginx-profile", "support-helm-urls")
	//			gitRepo := sourcev1.GitRepository{
	//				ObjectMeta: metav1.ObjectMeta{
	//					Name:      gitRefName,
	//					Namespace: namespace,
	//				},
	//				TypeMeta: metav1.TypeMeta{
	//					Kind:       "GitRepository",
	//					APIVersion: "source.toolkit.fluxcd.io/v1beta1",
	//				},
	//				Spec: sourcev1.GitRepositorySpec{
	//					URL: profileURL,
	//					Reference: &sourcev1.GitRepositoryRef{
	//						Branch: "support-helm-urls",
	//					},
	//				},
	//			}
	//			Expect(k8sClient.Create(ctx, &gitRepo)).Should(Succeed())
	//
	//			pSub := profilesv1.ProfileSubscription{
	//				TypeMeta: metav1.TypeMeta{
	//					Kind:       "ProfileSubscription",
	//					APIVersion: "profilesubscriptions.weave.works/v1alpha1",
	//				},
	//				ObjectMeta: metav1.ObjectMeta{
	//					Name:      subscriptionName,
	//					Namespace: namespace,
	//				},
	//				Spec: profilesv1.ProfileSubscriptionSpec{
	//					ProfileURL: profileURL,
	//				},
	//			}
	//			Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())
	//
	//			profile := profilesv1.ProfileSubscription{}
	//			Eventually(func() bool {
	//				err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
	//				return err == nil && profile.Status != profilesv1.ProfileSubscriptionStatus{}
	//			}, 10*time.Second, 1*time.Second).Should(BeTrue())
	//
	//			Expect(profile.Status.Message).To(Equal("error when creating profile artifacts"))
	//			Expect(profile.Status.State).To(Equal("failing"))
	//		})
	//	})
	//})
	Context("Create Remote", func() {
		DescribeTable("Applying a Profile creates the correct resources", func(pSubSpec profilesv1.ProfileSubscriptionSpec) {
			subscriptionName := "foo"
			branch := pSubSpec.Branch

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

			By("creating a HelmRepository resource")
			profileRepoName := "nginx-profile"
			gitRepoName := fmt.Sprintf("%s-%s-%s-remote", subscriptionName, profileRepoName, branch)
			helmRepository := sourcev1.HelmRepository{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: gitRepoName, Namespace: namespace}, &helmRepository)
			}, 10*time.Second).ShouldNot(HaveOccurred())
			Expect(helmRepository.Spec.URL).To(Equal("https://charts.bitnami.com/bitnami"))

			By("creating a HelmRelease resource")
			profileName := "nginx"
			chartName := "nginx-server"
			helmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName)
			helmRelease := helmv2.HelmRelease{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
			}, 10*time.Second).ShouldNot(HaveOccurred())
			Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal("nginx"))
			Expect(helmRelease.Spec.Chart.Spec.Version).To(Equal("8.8.1"))
			Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
				helmv2.CrossNamespaceObjectReference{
					Kind:      "HelmRepository",
					Name:      gitRepoName,
					Namespace: namespace,
				},
			))
			if pSub.Spec.Values != nil {
				Expect(helmRelease.Spec.Values).To(Equal(pSub.Spec.Values))
			}
			if pSub.Spec.ValuesFrom != nil {
				Expect(helmRelease.Spec.ValuesFrom).To(Equal(pSub.Spec.ValuesFrom))
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

			By("updating the status to Ready Unknown if the artifact resource reports it")
			profile := profilesv1.ProfileSubscription{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)).To(Succeed())
				return len(profile.Status.Conditions) > 0 &&
					profile.Status.Conditions[0].Type == "Ready" &&
					profile.Status.Conditions[0].Status == metav1.ConditionStatus("Unknown") &&
					profile.Status.Conditions[0].Message == "somethings wrong"
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

			helmResNew := helmRelease.DeepCopyObject().(*helmv2.HelmRelease)
			helmResNew.Status.Conditions = conditions
			Expect(k8sClient.Status().Patch(ctx, helmResNew, client.MergeFrom(&helmRelease))).To(Succeed())

			By("updating the status to Ready True when the resources are reporting Ready")
			profile = profilesv1.ProfileSubscription{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)).To(Succeed())
				return len(profile.Status.Conditions) > 0 &&
					profile.Status.Conditions[0].Type == "Ready" &&
					profile.Status.Conditions[0].Status == metav1.ConditionStatus("True") &&
					profile.Status.Conditions[0].Message == "all artifact resouces ready"
			}, 10*time.Second).Should(BeTrue())
		},
			Entry("a single Helm chart with no supplied values", profilesv1.ProfileSubscriptionSpec{
				ProfileURL: nginxProfileURL,
				Branch:     "support-helm-urls",
			}),
			Entry("a single Helm chart with supplied values", profilesv1.ProfileSubscriptionSpec{
				ProfileURL: nginxProfileURL,
				Branch:     "support-helm-urls",
				Values: &apiextensionsv1.JSON{
					Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
				},
			}),
			Entry("a single Helm chart with values supplied via valuesFrom", profilesv1.ProfileSubscriptionSpec{
				ProfileURL: nginxProfileURL,
				Branch:     "support-helm-urls",
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
				helmRepositoryUrl := "https://charts.bitnami.com/bitnami"

				gitRefName := fmt.Sprintf("%s-%s-%s-remote", subscriptionName, "nginx-profile", "support-helm-urls")
				helmRepository := sourcev1.HelmRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      gitRefName,
						Namespace: namespace,
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "HelmRepository",
						APIVersion: "source.toolkit.fluxcd.io/v1beta1",
					},
					Spec: sourcev1.HelmRepositorySpec{
						URL: helmRepositoryUrl,
					},
				}
				Expect(k8sClient.Create(ctx, &helmRepository)).Should(Succeed())

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
						ProfileURL: nginxProfileURL,
						Branch:     "support-helm-urls",
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				profile := profilesv1.ProfileSubscription{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: subscriptionName, Namespace: namespace}, &profile)
					return err == nil && len(profile.Status.Conditions) > 0
				}, 10*time.Second, 1*time.Second).Should(BeTrue())

				Expect(profile.Status.Conditions[0].Message).To(Equal("error when creating profile artifacts"))
				Expect(profile.Status.Conditions[0].Type).To(Equal("Ready"))
				Expect(profile.Status.Conditions[0].Status).To(Equal(metav1.ConditionStatus("False")))
				Expect(profile.Status.Conditions[0].Reason).To(Equal("CreateFailed"))
			})
		})
	})
})
