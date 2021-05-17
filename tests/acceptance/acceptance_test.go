package acceptance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

const (
	profileSubscriptionKind       = "ProfileSubscription"
	profileSubscriptionAPIVersion = "profilesubscriptions.weave.works/v1alpha1"

	nginxImage = "docker.io/bitnami/nginx:1.19.8-debian-10-r0"
)

var _ = Describe("Acceptance", func() {
	Context("ProfileSubscription", func() {
		var (
			profileURL string
			namespace  string
			branch     string
			subName    = "foo"
			nsp        v1.Namespace
		)

		BeforeEach(func() {
			profileURL = "https://github.com/weaveworks/profiles-examples"
			branch = "main"

			namespace = uuid.New().String()
			nsp = v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(kClient.Create(context.Background(), &nsp)).To(Succeed())
		})

		AfterEach(func() {
			Expect(kClient.Delete(context.Background(), &nsp)).To(Succeed())
		})

		When("subscribing to a Profile with a Helm Chart", func() {
			It("should deploy the Profile workload, reconcile when changes occur and cleanup on deletion", func() {
				pSub := profilesv1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       profileSubscriptionKind,
						APIVersion: profileSubscriptionAPIVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
						Version:    "weaveworks-nginx/v0.1.0",
					},
				}
				Expect(kClient.Create(context.Background(), &pSub)).To(Succeed())

				By("successfully deploying the helm release")
				helmReleaseName := fmt.Sprintf("%s-%s-%s", subName, "bitnami-nginx", "nginx-server")
				var helmRelease *helmv2.HelmRelease
				Eventually(func() bool {
					helmRelease = &helmv2.HelmRelease{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, helmRelease)
					if err != nil {
						return false
					}
					for _, condition := range helmRelease.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							return true
						}
					}
					return false
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				helmOpts := []client.ListOption{
					client.InNamespace(namespace),
					client.MatchingLabels{"app.kubernetes.io/name": "nginx"},
				}
				var podList *v1.PodList
				Eventually(func() v1.PodPhase {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, helmOpts...)
					Expect(err).NotTo(HaveOccurred())
					if len(podList.Items) == 0 {
						return v1.PodPhase("")
					}
					return podList.Items[0].Status.Phase
				}, 2*time.Minute, 5*time.Second).Should(Equal(v1.PodPhase("Running")))

				Expect(podList.Items[0].Spec.Containers[0].Image).To(Equal(nginxImage))

				By("successfully deploying the kustomize resource")
				kustomizeName := fmt.Sprintf("%s-%s-%s", subName, "weaveworks-nginx", "nginx-deployment")
				var kustomize *kustomizev1.Kustomization
				Eventually(func() bool {
					kustomize = &kustomizev1.Kustomization{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: kustomizeName, Namespace: namespace}, kustomize)
					if err != nil {
						return false
					}
					for _, condition := range kustomize.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							return true
						}
					}
					return false
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				kustomizeOpts := []client.ListOption{
					client.InNamespace(namespace),
					client.MatchingLabels{"app": "nginx"},
				}
				Eventually(func() v1.PodPhase {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, kustomizeOpts...)
					Expect(err).NotTo(HaveOccurred())
					if len(podList.Items) == 0 {
						return v1.PodPhase("no pods found")
					}
					return podList.Items[0].Status.Phase
				}, 2*time.Minute, 5*time.Second).Should(Equal(v1.PodPhase("Running")))

				Expect(podList.Items[0].Spec.Containers[0].Image).To(Equal("nginx:1.14.2"))

				By("recreating deleted artifacts")
				kustomize = &kustomizev1.Kustomization{}
				err := kClient.Get(context.Background(), client.ObjectKey{Name: kustomizeName, Namespace: namespace}, kustomize)
				Expect(err).NotTo(HaveOccurred())
				err = kClient.Delete(context.Background(), kustomize)
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() bool {
					kustomize = &kustomizev1.Kustomization{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: kustomizeName, Namespace: namespace}, kustomize)
					if err != nil {
						return false
					}
					for _, condition := range kustomize.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							return true
						}
					}
					return false
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				Eventually(func() v1.PodPhase {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, kustomizeOpts...)
					Expect(err).NotTo(HaveOccurred())
					if len(podList.Items) == 0 {
						return v1.PodPhase("no pods found")
					}
					return podList.Items[0].Status.Phase
				}, 2*time.Minute, 5*time.Second).Should(Equal(v1.PodPhase("Running")))

				Expect(podList.Items[0].Spec.Containers[0].Image).To(Equal("nginx:1.14.2"))

				By("cleaning up resources on deletion")
				Expect(kClient.Delete(context.Background(), &pSub)).To(Succeed())

				Eventually(func() bool {
					kustomize = &kustomizev1.Kustomization{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: kustomizeName, Namespace: namespace}, kustomize)
					return apierrors.IsNotFound(err)
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				Eventually(func() int {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, kustomizeOpts...)
					Expect(err).NotTo(HaveOccurred())
					return len(podList.Items)
				}, 5*time.Minute, 10*time.Second).Should(Equal(0))

				Eventually(func() bool {
					helmRelease = &helmv2.HelmRelease{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, helmRelease)
					return apierrors.IsNotFound(err)
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				Eventually(func() int {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, helmOpts...)
					Expect(err).NotTo(HaveOccurred())
					return len(podList.Items)
				}, 5*time.Minute, 10*time.Second).Should(Equal(0))

			})
		})
		When("subscribing to a Profile with a Helm Chart using chart description", func() {
			It("should deploy the Profile workload", func() {
				pCatalog := profilesv1.ProfileCatalogSource{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileCatalogSource",
						APIVersion: profileSubscriptionAPIVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "weaveworks-nginx",
						Namespace: "default",
					},
					Spec: profilesv1.ProfileCatalogSourceSpec{
						Profiles: []profilesv1.ProfileDescription{
							{
								Name:    "bitnami-nginx",
								Version: "v0.1.0",
								URL:     profileURL,
							},
						},
					},
				}
				Expect(kClient.Create(context.Background(), &pCatalog)).To(Succeed())

				pSub := profilesv1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       profileSubscriptionKind,
						APIVersion: profileSubscriptionAPIVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileSubscriptionSpec{
						ProfileCatalogDescription: &profilesv1.ProfileCatalogDescription{
							Version: "v0.1.0",
							Catalog: "weaveworks-nginx",
							Profile: "bitnami-nginx",
						},
					},
				}
				Expect(kClient.Create(context.Background(), &pSub)).To(Succeed())

				By("successfully deploying the helm release")
				helmReleaseName := fmt.Sprintf("%s-%s-%s", subName, "bitnami-nginx", "nginx-server")
				var helmRelease *helmv2.HelmRelease
				Eventually(func() bool {
					helmRelease = &helmv2.HelmRelease{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, helmRelease)
					if err != nil {
						return false
					}
					for _, condition := range helmRelease.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							return true
						}
					}
					return false
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())
			})
		})
		When("updating values in a subscription", func() {
			It("should reconcile and apply the new values", func() {
				pSub := profilesv1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       profileSubscriptionKind,
						APIVersion: profileSubscriptionAPIVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
						Branch:     branch,
						Path:       "weaveworks-nginx",
						Values: &apiextensionsv1.JSON{
							Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
						},
					},
				}
				Expect(kClient.Create(context.Background(), &pSub)).To(Succeed())
				appName := "bitnami-nginx"
				By("successfully deploying the helm release")
				helmReleaseName := fmt.Sprintf("%s-%s-%s", subName, appName, "nginx-server")
				var helmRelease *helmv2.HelmRelease
				Eventually(func() bool {
					helmRelease = &helmv2.HelmRelease{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, helmRelease)
					if err != nil {
						return false
					}
					for _, condition := range helmRelease.Status.Conditions {
						if condition.Type == "Ready" && condition.Status == "True" {
							return true
						}
					}
					return false
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				By("updating the values")
				updatedSub := &profilesv1.ProfileSubscription{}
				Eventually(func() error {
					return kClient.Get(context.Background(), client.ObjectKey{Name: subName, Namespace: namespace}, updatedSub)
				}, 2*time.Minute, 5*time.Second).ShouldNot(HaveOccurred())
				updatedRaw := []byte(`{"replicaCount":1,"service":{"port":8881}}`)
				updatedSub.Spec.Values.Raw = updatedRaw
				Expect(kClient.Patch(context.Background(), updatedSub, client.MergeFrom(&pSub))).To(Succeed())

				helmReleaseNew := helmv2.HelmRelease{}
				Eventually(func() bool {
					err := kClient.Get(context.Background(), client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmReleaseNew)
					return err == nil && bytes.Equal(helmReleaseNew.Spec.Values.Raw, updatedRaw)
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				By("examining the pod count as we change replica count value")
				Eventually(func() bool {
					pods := &v1.PodList{}
					labelSelector, _ := labels.Parse("app.kubernetes.io/instance=" + helmReleaseName)
					opts := &client.ListOptions{
						Namespace:     namespace,
						LabelSelector: labelSelector,
					}
					err := kClient.List(context.Background(), pods, opts)
					return err == nil && len(pods.Items) == 1
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())
			})
		})
	})

	Context("ProfileCatalog", func() {
		var (
			pCatalog                profilesv1.ProfileCatalogSource
			expectedNginx1          profilesv1.ProfileDescription
			sourceName, profileName string
		)

		BeforeEach(func() {
			sourceName, profileName = "catalog", "nginx-1"
			pCatalog = profilesv1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: profileSubscriptionAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      sourceName,
					Namespace: "default",
				},
				Spec: profilesv1.ProfileCatalogSourceSpec{
					Profiles: []profilesv1.ProfileDescription{
						{
							Name:          profileName,
							Description:   "nginx 1",
							Version:       "0.0.1",
							URL:           "foo.com/bar",
							Maintainer:    "my aunt ethel",
							Prerequisites: []string{"at least 20 years of kubernetes experience"},
						},
						{
							Name:        "nginx-2",
							Description: "nginx 1",
						},
						{
							Name:        "something-else",
							Description: "something else",
						},
					},
				},
			}
			Expect(kClient.Create(context.Background(), &pCatalog)).To(Succeed())

			expectedNginx1 = profilesv1.ProfileDescription{
				Name:          profileName,
				Description:   "nginx 1",
				CatalogSource: sourceName,
				Version:       "0.0.1",
				URL:           "foo.com/bar",
				Maintainer:    "my aunt ethel",
				Prerequisites: []string{"at least 20 years of kubernetes experience"},
			}
		})

		AfterEach(func() {
			_ = kClient.Delete(context.Background(), &pCatalog)
		})

		Context("search", func() {
			It("returns the matching catalogs", func() {
				Eventually(func() []profilesv1.ProfileDescription {
					req, err := http.NewRequest("GET", "http://localhost:8000/profiles", nil)
					Expect(err).NotTo(HaveOccurred())
					u, err := url.Parse("http://localhost:8000")
					Expect(err).NotTo(HaveOccurred())
					q := u.Query()
					q.Add("name", "nginx")
					req.URL.RawQuery = q.Encode()
					Expect(err).NotTo(HaveOccurred())
					resp, err := http.DefaultClient.Do(req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					descriptions := []profilesv1.ProfileDescription{}
					_ = json.NewDecoder(resp.Body).Decode(&descriptions)
					return descriptions
				}).Should(ConsistOf(
					expectedNginx1,
					profilesv1.ProfileDescription{
						Name:          "nginx-2",
						Description:   "nginx 1",
						CatalogSource: sourceName,
					},
				))
			})
		})

		Context("get", func() {
			It("returns details of the requested catalog entry", func() {
				Eventually(func() profilesv1.ProfileDescription {
					description, _ := getProfile(profileName, sourceName)
					return description
				}, "10s").Should(Equal(expectedNginx1))
			})
		})

		Context("update", func() {
			It("updates a ProfileCatalogSource with new profiles", func() {
				pCatalog.Spec.Profiles = append(pCatalog.Spec.Profiles, profilesv1.ProfileDescription{
					Name:        "new-profile",
					Description: "I am new here",
				})
				Expect(kClient.Update(context.Background(), &pCatalog)).To(Succeed())
				Eventually(func() profilesv1.ProfileDescription {
					description, err := getProfile("new-profile", sourceName)
					Expect(err).NotTo(HaveOccurred())
					return description
				}).Should(Equal(profilesv1.ProfileDescription{
					Name:          "new-profile",
					Description:   "I am new here",
					CatalogSource: sourceName,
				}))
			})
		})

		Context("delete", func() {
			It("clears the in-memory cache when a ProfileCatalogSource is deleted", func() {
				description, err := getProfile(profileName, sourceName)
				Expect(err).NotTo(HaveOccurred())
				Expect(description).To(Equal(expectedNginx1))

				Expect(kClient.Delete(context.Background(), &pCatalog)).To(Succeed())
				Eventually(func() error {
					_, err := getProfile(profileName, sourceName)
					return err
				}, "5s").Should(MatchError(ContainSubstring("got 404")))
			})
		})
	})
})

func getProfile(profileName, sourceName string) (profilesv1.ProfileDescription, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/profiles/%s/%s", sourceName, profileName))
	if err != nil {
		return profilesv1.ProfileDescription{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return profilesv1.ProfileDescription{}, fmt.Errorf("expected status code 200; got %d", resp.StatusCode)
	}
	var p profilesv1.ProfileDescription
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return profilesv1.ProfileDescription{}, err
	}
	return p, nil
}
