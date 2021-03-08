package integration_test

import (
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	binaryPath string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PCTL Integration Suite")
}

var _ = BeforeSuite(func() {
	var err error
	binaryPath, err = gexec.Build("github.com/weaveworks/profiles/cmd/pctl")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
