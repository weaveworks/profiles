package catalog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Catalog", func() {
	Context("Search", func() {
		It("returns matching profiles", func() {
			c := catalog.New()
			c.Add(profilesv1.ProfileDescription{Name: "foo"})
			c.Add(profilesv1.ProfileDescription{Name: "bar"})
			c.Add(profilesv1.ProfileDescription{Name: "alsofoo"})
			Expect(c.Search("foo")).To(ConsistOf(
				profilesv1.ProfileDescription{Name: "foo"},
				profilesv1.ProfileDescription{Name: "alsofoo"},
			))
		})
	})
})
