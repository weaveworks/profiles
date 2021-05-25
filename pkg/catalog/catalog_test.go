package catalog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
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
		profiles := []profilesv1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
		c.Update(catName, profiles...)

		By("returning all matching profiles based on query string")
		Expect(c.Search("foo")).To(ConsistOf(
			profilesv1.ProfileDescription{Name: "foo", CatalogSource: catName},
			profilesv1.ProfileDescription{Name: "alsofoo", CatalogSource: catName},
		))

		By("getting details for a specific named profile in a catalog")
		Expect(c.Get(catName, "foo")).To(Equal(
			&profilesv1.ProfileDescription{Name: "foo", CatalogSource: catName},
		))

		By("updating profiles in a catalog source")
		profiles = []profilesv1.ProfileDescription{{Name: "foo"}, {Name: "bar"}}
		c.Update(catName, profiles...)
		Expect(c.Search("foo")).To(ConsistOf(profilesv1.ProfileDescription{Name: "foo", CatalogSource: catName}))

		By("removing a catalog source")
		c.Remove(catName)
		Expect(c.Search("foo")).To(BeEmpty())
	})

	Describe("GetWithVersion", func() {
		It("returns the profile with the matching verrsion", func() {

			profiles := []profilesv1.ProfileDescription{{Name: "foo", Version: "v0.1.0", Description: "install foo"}, {Name: "foo", Version: "0.2.0"}, {Name: "foo"}}
			c.Update(catName, profiles...)

			Expect(c.GetWithVersion(catName, "foo", "v0.1.0")).To(Equal(
				&profilesv1.ProfileDescription{Name: "foo", Version: "v0.1.0", Description: "install foo", CatalogSource: catName},
			))
		})

		When("version is set to latest", func() {
			It("returns the latest version", func() {
				profiles := []profilesv1.ProfileDescription{{Name: "foo", Version: "v0.1.0"}, {Name: "foo", Version: "0.2.0"}, {Name: "foo"}}
				c.Update(catName, profiles...)

				Expect(c.GetWithVersion(catName, "foo", "latest")).To(Equal(
					&profilesv1.ProfileDescription{Name: "foo", Version: "0.2.0", CatalogSource: catName},
				))

				profiles = []profilesv1.ProfileDescription{{Name: "foo", Version: "v0.3.0"}, {Name: "foo", Version: "0.2.0"}, {Name: "foo"}}
				c.Update(catName, profiles...)

				Expect(c.GetWithVersion(catName, "foo", "latest")).To(Equal(
					&profilesv1.ProfileDescription{Name: "foo", Version: "v0.3.0", CatalogSource: catName},
				))
			})

			When("no profile has a valid version", func() {
				It("returns nil", func() {
					profiles := []profilesv1.ProfileDescription{{Name: "foo", Version: "vsda012!.1.0"}, {Name: "foo", Version: "!0.!2.0"}, {Name: "foo"}}
					c.Update(catName, profiles...)

					Expect(c.GetWithVersion(catName, "foo", "latest")).To(BeNil())
				})
			})
		})
	})
})
