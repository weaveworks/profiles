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
		log logr.Logger
		c   *catalog.Catalog
	)

	BeforeEach(func() {
		log = logr.Discard()
		c = catalog.New(log)
	})

	Context("Add", func() {
		It("does not add duplicate profiles", func() {
			// 2 foos added
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "foo"}, {Name: "bar"}}
			c.Add(profiles...)
			// only 1 should remain
			Expect(c.Search("foo")).To(ConsistOf(v1alpha1.ProfileDescription{Name: "foo"}))
		})
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
