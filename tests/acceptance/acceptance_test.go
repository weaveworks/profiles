package acceptance_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	profileSubscriptionKind       = "ProfileSubscription"
	profileSubscriptionAPIVersion = "profilesubscriptions.weave.works/v1alpha1"

	nginxImage = "docker.io/bitnami/nginx:1.19.8-debian-10-r0"
)

var _ = Describe("Acceptance", func() {
	var (
		profileURL, namespace string

		nsp v1.Namespace
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
		It("should deploy the Profile workload", func() {
			pSub := v1alpha1.ProfileSubscription{
				TypeMeta: metav1.TypeMeta{
					Kind:       profileSubscriptionKind,
					APIVersion: profileSubscriptionAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: namespace,
				},
				Spec: v1alpha1.ProfileSubscriptionSpec{
					ProfileURL: profileURL,
				},
			}
			Expect(kClient.Create(context.Background(), &pSub)).To(Succeed())

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
			}, 2*time.Minute, 10*time.Second).Should(Equal(v1.PodPhase("Running")))

			// TODO we check the Pod because for some reason we the HelmRelease will
			// not register as running, even though all child resources are fine
			// :shrug_emoji:
			Expect(podList.Items[0].Spec.Containers[0].Image).To(Equal(nginxImage))
		})
	})
})
