package acceptance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	profileSubscriptionKind       = "ProfileSubscription"
	profileSubscriptionAPIVersion = "profilesubscriptions.weave.works/v1alpha1"

	nginxImage = "docker.io/bitnami/nginx:1.19.8-debian-10-r0"
)

var _ = Describe("Acceptance", func() {
	Context("ProfileSubscription", func() {
		var (
			profileURL, namespace string
			subName               = "foo"
			nsp                   v1.Namespace
		)

		BeforeEach(func() {
			profileURL = "https://github.com/weaveworks/nginx-profile"

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
				pSub := v1alpha1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       profileSubscriptionKind,
						APIVersion: profileSubscriptionAPIVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      subName,
						Namespace: namespace,
					},
					Spec: v1alpha1.ProfileSubscriptionSpec{
						ProfileURL: profileURL,
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

				opts := []client.ListOption{
					client.InNamespace(namespace),
					client.MatchingLabels{"app.kubernetes.io/name": "nginx"},
				}
				var podList *v1.PodList
				Eventually(func() v1.PodPhase {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, opts...)
					Expect(err).NotTo(HaveOccurred())
					if len(podList.Items) == 0 {
						return v1.PodPhase("")
					}
					return podList.Items[0].Status.Phase
				}, 2*time.Minute, 5*time.Second).Should(Equal(v1.PodPhase("Running")))

				Expect(podList.Items[0].Spec.Containers[0].Image).To(Equal(nginxImage))

				By("cleaning up resources on deletion")
				Expect(kClient.Delete(context.Background(), &pSub)).To(Succeed())

				Eventually(func() bool {
					helmRelease = &helmv2.HelmRelease{}
					err := kClient.Get(context.Background(), client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, helmRelease)
					return apierrors.IsNotFound(err)
				}, 2*time.Minute, 5*time.Second).Should(BeTrue())

				Eventually(func() int {
					podList = &v1.PodList{}
					err := kClient.List(context.Background(), podList, opts...)
					Expect(err).NotTo(HaveOccurred())
					return len(podList.Items)
				}, 5*time.Minute, 10*time.Second).Should(Equal(0))

			})
		})
	})

	Context("ProfileCatalog", func() {
		It("returns the matching catalogs", func() {
			pCatalog := v1alpha1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: profileSubscriptionAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "catalog",
					Namespace: "default",
				},
				Spec: v1alpha1.ProfileCatalogSourceSpec{
					Profiles: []v1alpha1.ProfileDescription{
						{
							Name:        "nginx-1",
							Description: "nginx 1",
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
			Eventually(func() []v1alpha1.ProfileDescription {
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
				descriptions := []v1alpha1.ProfileDescription{}
				_ = json.NewDecoder(resp.Body).Decode(&descriptions)
				return descriptions

			}).Should(ConsistOf(
				v1alpha1.ProfileDescription{
					Name:        "nginx-1",
					Description: "nginx 1",
				},
				v1alpha1.ProfileDescription{
					Name:        "nginx-2",
					Description: "nginx 1",
				},
			))

		})
	})
})
