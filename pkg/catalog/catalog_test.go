package catalog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Catalog", func() {
	Context("Search", func() {
		It("returns matching profiles", func() {
			c := catalog.New()
			c.Add(v1alpha1.ProfileDescription{Name: "foo"})
			c.Add(v1alpha1.ProfileDescription{Name: "bar"})
			c.Add(v1alpha1.ProfileDescription{Name: "alsofoo"})
			Expect(c.Search("foo")).To(ConsistOf(
				v1alpha1.ProfileDescription{Name: "foo"},
				v1alpha1.ProfileDescription{Name: "alsofoo"},
			))
		})
	})
})
