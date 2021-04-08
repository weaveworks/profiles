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

	Context("Add", func() {
		BeforeEach(func() {
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(catName, profiles...)
		})

		It("adds the catalog name to each profile", func() {
			Expect(c.Get("foo")).To(Equal(
				v1alpha1.ProfileDescription{Name: "foo", Catalog: catName},
			))
		})
	})

	Context("Search", func() {
		BeforeEach(func() {
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(catName, profiles...)
		})

		It("returns matching profiles", func() {
			Expect(c.Search("foo")).To(ConsistOf(
				v1alpha1.ProfileDescription{Name: "foo", Catalog: catName},
				v1alpha1.ProfileDescription{Name: "alsofoo", Catalog: catName},
			))
		})
	})

	Context("Get", func() {
		BeforeEach(func() {
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(catName, profiles...)
		})

		It("returns the requested profile", func() {
			Expect(c.Get("foo")).To(Equal(
				v1alpha1.ProfileDescription{Name: "foo", Catalog: catName},
			))
		})
	})
})
