package version_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/pkg/version"
)

var _ = Describe("Version", func() {
	Context("Get", func() {
		It("should return the version", func() {
			Expect(version.Get()).To(Equal("v0.0.1"))
		})
	})
})
