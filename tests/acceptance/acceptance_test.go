package acceptance_test

import (
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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			profileURL = "https://github.com/weaveworks/nginx-profile"
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
			It("should deploy the Profile workload and cleanup on deletion", func() {
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
					},
				}
				Expect(kClient.Create(context.Background(), &pSub)).To(Succeed())

				By("successfully deploying the helm release")
				helmReleaseName := fmt.Sprintf("%s-%s-%s", subName, "nginx", "nginx-server")
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
				kustomizeName := fmt.Sprintf("%s-%s-%s", subName, "nginx", "nginx-deployment")
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
	})

	Context("ProfileCatalog", func() {
		var (
			pCatalog                 profilesv1.ProfileCatalogSource
			expectedNginx1           profilesv1.ProfileDescription
			catalogName, profileName string
		)

		BeforeEach(func() {
			catalogName, profileName = "catalog", "nginx-1"
			pCatalog = profilesv1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: profileSubscriptionAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      catalogName,
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
				Catalog:       catalogName,
				Version:       "0.0.1",
				URL:           "foo.com/bar",
				Maintainer:    "my aunt ethel",
				Prerequisites: []string{"at least 20 years of kubernetes experience"},
			}
		})

		AfterEach(func() {
			Expect(kClient.Delete(context.Background(), &pCatalog)).To(Succeed())
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
					Expect(json.NewDecoder(resp.Body).Decode(&descriptions)).To(Succeed())
					return descriptions
				}).Should(ConsistOf(
					expectedNginx1,
					profilesv1.ProfileDescription{
						Name:        "nginx-2",
						Description: "nginx 1",
						Catalog:     catalogName,
					},
				))
			})
		})

		makeProfileURL := func() string {
			return fmt.Sprintf("http://localhost:8000/profiles/%s/%s", catalogName, profileName)
		}

		getProfileDescription := func() (profilesv1.ProfileDescription, error) {
			resp, err := http.Get(makeProfileURL())
			if err != nil {
				return profilesv1.ProfileDescription{}, err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return profilesv1.ProfileDescription{}, fmt.Errorf("expected status code 200; got %d", resp.StatusCode)
			}
			var description profilesv1.ProfileDescription
			if err := json.NewDecoder(resp.Body).Decode(&description); err != nil {
				return profilesv1.ProfileDescription{}, err
			}
			return description, nil
		}

		Context("get and delete", func() {
			It("returns details of the requested catalog entry and deletes it", func() {
				By("fetching the catalog entry")
				Eventually(func() profilesv1.ProfileDescription {
					description, err := getProfileDescription()
					Expect(err).NotTo(HaveOccurred())
					return description
				}).Should(Equal(expectedNginx1))

				By("deleting the catalog")
				client := http.Client{
					Timeout: 5 * time.Second,
				}
				req, err := http.NewRequest(http.MethodDelete, makeProfileURL(), nil)
				Expect(err).NotTo(HaveOccurred())
				resp, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				resp, err = client.Get(makeProfileURL())
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

	})
})
