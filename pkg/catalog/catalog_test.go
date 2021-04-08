package catalog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Catalog", func() {
	var (
		c       *catalog.Catalog
		catName string
	)

	BeforeEach(func() {
		c = catalog.New()
		catName = "whiskers"
	})

	It("manages an in memory list of profiles", func() {
		By("adding profiles to the list")
		profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
		c.Add(catName, profiles...)

		By("returning all matching profiles based on query string")
		Expect(c.Search("foo")).To(ConsistOf(
			v1alpha1.ProfileDescription{Name: "foo", Catalog: catName},
			v1alpha1.ProfileDescription{Name: "alsofoo", Catalog: catName},
		))

		By("getting details for a specific named profile in a catalog")
		Expect(c.Get(catName, "foo")).To(Equal(
			v1alpha1.ProfileDescription{Name: "foo", Catalog: catName},
		))
	})
})
