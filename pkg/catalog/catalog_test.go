package catalog_test

import (
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Catalog", func() {
	var log = logr.Discard()
	Context("Add", func() {
		It("adds profiles to the catalog", func() {
			c := catalog.New(log)
			c.Add(v1alpha1.ProfileDescription{Name: "foo"})
			Expect(c.Profiles()).To(Equal(
				[]v1alpha1.ProfileDescription{{Name: "foo"}},
			))
		})

		It("does not add duplicate profiles", func() {
			c := catalog.New(log)
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "foo"}}
			c.Add(profiles...)
			Expect(c.Profiles()).To(HaveLen(2))
			Expect(c.Profiles()).To(Equal(
				[]v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}},
			))
		})
	})

	Context("Search", func() {
		It("returns matching profiles", func() {
			c := catalog.New(log)
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(profiles...)
			Expect(c.Search("foo")).To(ConsistOf(
				profilesv1.ProfileDescription{Name: "foo"},
				profilesv1.ProfileDescription{Name: "alsofoo"},
			))
		})
	})

	Context("Show", func() {
		It("returns the requested profile", func() {
			c := catalog.New(log)
			profiles := []v1alpha1.ProfileDescription{{Name: "foo"}, {Name: "bar"}, {Name: "alsofoo"}}
			c.Add(profiles...)
			Expect(c.Show("foo")).To(Equal(
				v1alpha1.ProfileDescription{Name: "foo"},
			))
		})
	})
})
