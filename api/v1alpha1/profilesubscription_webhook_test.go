package v1alpha1_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ProfilesubscriptionWebhook", func() {
	Context("defaulting", func() {
		When("the branch is not set", func() {
			It("defaults to main", func() {
				pSub := v1alpha1.ProfileSubscription{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileSubscription",
						APIVersion: "profilesubscriptions.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: v1alpha1.ProfileSubscriptionSpec{
						ProfileURL: "github.com/foo/bar",
					},
				}
				Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

				Eventually(func() v1alpha1.ProfileSubscriptionSpec {
					pSub = v1alpha1.ProfileSubscription{}
					err := k8sClient.Get(ctx, client.ObjectKey{Name: "test", Namespace: "default"}, &pSub)
					Expect(err).NotTo(HaveOccurred())
					return pSub.Spec
				}, 5*time.Second).Should(Equal(v1alpha1.ProfileSubscriptionSpec{
					ProfileURL: "github.com/foo/bar",
					Branch:     "main",
				}))
			})
		})
	})
})
