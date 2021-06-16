package catalog_test

import (
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Catalog", func() {
	var (
		c       *catalog.Catalog
		catName string
		logger  = logr.Discard()
	)

	BeforeEach(func() {
		c = catalog.New()
		catName = "whiskers"
	})

	It("manages an in memory list of profiles", func() {
		By("adding profiles to the list")
		profiles := []profilesv1.ProfileCatalogEntry{
			{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}},
			{ProfileDescription: profilesv1.ProfileDescription{Name: "bar"}},
			{ProfileDescription: profilesv1.ProfileDescription{Name: "alsofoo"}},
		}
		c.Update(catName, profiles...)

		By("returning all the profiles available")
		Expect(c.SearchAll()).To(ConsistOf(
			profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "bar"}, CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "alsofoo"}, CatalogSource: catName},
		))

		By("returning all matching profiles based on query string")
		Expect(c.Search("foo")).To(ConsistOf(
			profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "alsofoo"}, CatalogSource: catName},
		))

		By("getting details for a specific named profile in a catalog")
		Expect(c.Get(catName, "foo")).To(Equal(
			&profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, CatalogSource: catName},
		))

		By("updating profiles in a catalog source")
		profiles = []profilesv1.ProfileCatalogEntry{
			{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}},
			{ProfileDescription: profilesv1.ProfileDescription{Name: "bar"}},
		}
		c.Update(catName, profiles...)
		Expect(c.Search("foo")).To(ConsistOf(profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, CatalogSource: catName}))

		By("removing a catalog source")
		c.Remove(catName)
		Expect(c.Search("foo")).To(BeEmpty())
	})

	Describe("GetWithVersion", func() {
		It("returns the profile with the matching version", func() {

			profiles := []profilesv1.ProfileCatalogEntry{
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo", Description: "install foo"}, Tag: "v0.1.0"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "0.2.0"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}}}
			c.Update(catName, profiles...)

			Expect(c.GetWithVersion(logger, catName, "foo", "v0.1.0")).To(Equal(
				&profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo", Description: "install foo"}, Tag: "v0.1.0", CatalogSource: catName},
			))
		})

		When("version is set to latest", func() {
			It("returns the latest version", func() {
				profiles := []profilesv1.ProfileCatalogEntry{{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "foo/v0.1.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "foo/0.2.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "bar"}, Tag: "bar/0.3.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}}}
				c.Update(catName, profiles...)

				Expect(c.GetWithVersion(logger, catName, "foo", "latest")).To(Equal(
					&profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "foo/0.2.0", CatalogSource: catName},
				))

				profiles = []profilesv1.ProfileCatalogEntry{{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "0.2.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "v0.3.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}}}
				c.Update(catName, profiles...)
				Expect(c.GetWithVersion(logger, catName, "foo", "latest")).To(Equal(
					&profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "v0.3.0", CatalogSource: catName},
				))
			})

			When("no profile has a valid version", func() {
				It("returns nil", func() {
					profiles := []profilesv1.ProfileCatalogEntry{{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "vsda012!.1.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "!0.!2.0"}, {ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}}}
					c.Update(catName, profiles...)

					Expect(c.GetWithVersion(logger, catName, "foo", "latest")).To(BeNil())
				})
			})
		})
	})

	Describe("ProfilesGreaterThanVersion", func() {
		It("lists all available versions which are greater than the current version in descending order", func() {

			profiles := []profilesv1.ProfileCatalogEntry{
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo", Description: "install foo"}, Tag: "v0.0.1"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo", Description: "install foo"}, Tag: "v0.1.0"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "v0.2.0"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "v0.3.0"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo2"}, Tag: "v0.3.0"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo2"}, Tag: "v0.3.1"},
				{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}},
			}
			c.Update(catName, profiles...)

			Expect(c.ProfilesGreaterThanVersion(logger, catName, "foo", "v0.1.0")).To(Equal(
				[]profilesv1.ProfileCatalogEntry{
					{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "v0.3.0", CatalogSource: catName},
					{ProfileDescription: profilesv1.ProfileDescription{Name: "foo"}, Tag: "v0.2.0", CatalogSource: catName},
				},
			))
		})
	})
})
