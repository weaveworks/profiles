package controllers_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	Context("Create", func() {
		It("adds the profile to the in-memory list", func() {
			pSub := v1alpha1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: "profilesubscriptions.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "catalog",
					Namespace: namespace,
				},
				Spec: v1alpha1.ProfileCatalogSourceSpec{
					Profiles: []v1alpha1.ProfileDescription{
						{
							Name:        "foo",
							Description: "bar",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, &pSub)).Should(Succeed())

			Eventually(func() []v1alpha1.ProfileDescription {
				return catalogReconciler.Profiles.Search("foo")
			}, 2*time.Second).Should(ConsistOf(v1alpha1.ProfileDescription{Name: "foo", Description: "bar", Catalog: "catalog"}))
		})
	})
})
