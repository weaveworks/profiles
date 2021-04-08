package catalog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Catalog", func() {
	var c *catalog.Catalog

	BeforeEach(func() {
		c = catalog.New()
	})

	Context("Search", func() {
		BeforeEach(func() {
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(profiles...)
		})

		It("returns matching profiles", func() {
			Expect(c.Search("foo")).To(ConsistOf(
				profilesv1.ProfileDescription{Name: "foo"},
				profilesv1.ProfileDescription{Name: "alsofoo"},
			))
		})
	})

	Context("Get", func() {
		BeforeEach(func() {
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(profiles...)
		})

		It("returns the requested profile", func() {
			Expect(c.Get("foo")).To(Equal(
				v1alpha1.ProfileDescription{Name: "foo"},
			))
		})
	})
})
