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
		Expect(c.CatalogExists(catName)).To(BeFalse())
		profiles := []profilesv1.ProfileCatalogEntry{
			{Name: "foo"},
			{Name: "bar"},
			{Name: "alsofoo"},
		}
		c.Append(catName, profiles...)
		Expect(c.CatalogExists(catName)).To(BeTrue())

		By("returning all the profiles available")
		Expect(c.SearchAll()).To(ConsistOf(
			profilesv1.ProfileCatalogEntry{Name: "foo", CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{Name: "bar", CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{Name: "alsofoo", CatalogSource: catName},
		))

		By("returning all matching profiles based on query string")
		Expect(c.Search("foo")).To(ConsistOf(
			profilesv1.ProfileCatalogEntry{Name: "foo", CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{Name: "alsofoo", CatalogSource: catName},
		))

		By("getting details for a specific named profile in a catalog")
		Expect(c.Get(catName, "foo")).To(Equal(
			&profilesv1.ProfileCatalogEntry{Name: "foo", CatalogSource: catName},
		))

		By("replacing profiles in a catalog source")
		profiles = []profilesv1.ProfileCatalogEntry{
			{Name: "foo"},
			{Name: "bar"},
		}
		c.AddOrReplace(catName, profiles...)
		Expect(c.Search("foo")).To(ConsistOf(profilesv1.ProfileCatalogEntry{Name: "foo", CatalogSource: catName}))

		By("appending profiles in a catalog source")
		profiles = []profilesv1.ProfileCatalogEntry{
			{Name: "bar-2"},
		}
		c.Append(catName, profiles...)
		Expect(c.Search("bar")).To(ConsistOf(
			profilesv1.ProfileCatalogEntry{Name: "bar", CatalogSource: catName},
			profilesv1.ProfileCatalogEntry{Name: "bar-2", CatalogSource: catName},
		))

		By("removing a catalog source")
		c.Remove(catName)
		Expect(c.Search("foo")).To(BeEmpty())
	})

	Describe("GetWithVersion", func() {
		It("returns the profile with the matching version", func() {

			profiles := []profilesv1.ProfileCatalogEntry{
				{ProfileDescription: profilesv1.ProfileDescription{Description: "install foo"}, Name: "foo", Tag: "v0.1.0"},
				{Name: "foo", Tag: "0.2.0"},
				{Name: "foo"},
			}
			c.AddOrReplace(catName, profiles...)

			Expect(c.GetWithVersion(logger, catName, "foo", "v0.1.0")).To(Equal(
				&profilesv1.ProfileCatalogEntry{ProfileDescription: profilesv1.ProfileDescription{Description: "install foo"}, Name: "foo", Tag: "v0.1.0", CatalogSource: catName},
			))
		})

		When("version is set to latest", func() {
			It("returns the latest version", func() {
				profiles := []profilesv1.ProfileCatalogEntry{{Name: "foo", Tag: "foo/v0.1.0"}, {Name: "foo", Tag: "foo/0.2.0"}, {Name: "bar", Tag: "bar/0.3.0"}, {Name: "foo"}}
				c.AddOrReplace(catName, profiles...)

				Expect(c.GetWithVersion(logger, catName, "foo", "latest")).To(Equal(
					&profilesv1.ProfileCatalogEntry{Name: "foo", Tag: "foo/0.2.0", CatalogSource: catName},
				))

				profiles = []profilesv1.ProfileCatalogEntry{{Name: "foo", Tag: "0.2.0"}, {Name: "foo", Tag: "v0.3.0"}, {Name: "foo"}}
				c.AddOrReplace(catName, profiles...)
				Expect(c.GetWithVersion(logger, catName, "foo", "latest")).To(Equal(
					&profilesv1.ProfileCatalogEntry{Name: "foo", Tag: "v0.3.0", CatalogSource: catName},
				))
			})

			When("no profile has a valid version", func() {
				It("returns nil", func() {
					profiles := []profilesv1.ProfileCatalogEntry{{Name: "foo", Tag: "vsda012!.1.0"}, {Name: "foo", Tag: "!0.!2.0"}, {Name: "foo"}}
					c.AddOrReplace(catName, profiles...)

					Expect(c.GetWithVersion(logger, catName, "foo", "latest")).To(BeNil())
				})
			})
		})
	})

	Describe("ProfilesGreaterThanVersion", func() {
		It("lists all available versions which are greater than the current version in descending order", func() {

			profiles := []profilesv1.ProfileCatalogEntry{
				{ProfileDescription: profilesv1.ProfileDescription{Description: "install foo"}, Name: "foo", Tag: "v0.0.1"},
				{ProfileDescription: profilesv1.ProfileDescription{Description: "install foo"}, Name: "foo", Tag: "v0.1.0"},
				{Name: "foo", Tag: "v0.2.0"},
				{Name: "foo"},
				{Name: "foo", Tag: "v0.3.0"},
				{Name: "foo"},
				{Name: "foo2", Tag: "v0.3.0"},
				{Name: "foo"},
				{Name: "foo2", Tag: "v0.3.1"},
				{Name: "foo"},
			}
			c.AddOrReplace(catName, profiles...)

			Expect(c.ProfilesGreaterThanVersion(logger, catName, "foo", "v0.1.0")).To(Equal(
				[]profilesv1.ProfileCatalogEntry{
					{Name: "foo", Tag: "v0.3.0", CatalogSource: catName},
					{Name: "foo", Tag: "v0.2.0", CatalogSource: catName},
				},
			))
		})
	})
})
