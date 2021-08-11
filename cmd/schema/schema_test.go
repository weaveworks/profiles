package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	schemapkg "github.com/weaveworks/schemer/schema"
)

var _ = Describe("Schema", func() {
	var (
		tmpDir string
		object string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		object = "ProfileDefinition"
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("writes the generated schema of the given object to the given location", func() {
		destFile := filepath.Join(tmpDir, "schema.json")
		session, err := runCmd(object, destFile)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session.Out, 20).Should(gbytes.Say("schema file generated for `ProfileDefinition`"))

		Eventually(destFile, 20).Should(BeAnExistingFile())

		contents, err := ioutil.ReadFile(destFile)
		Expect(err).NotTo(HaveOccurred())
		s := schemapkg.Schema{}
		Expect(json.Unmarshal([]byte(contents), &s)).To(Succeed())
		Expect(s.Version).To(Equal("http://json-schema.org/draft-07/schema#"))
		Expect(s.Definition).ToNot(BeNil())
		Expect(s.Definitions).To(HaveLen(9))
		// I considered checking each thing here, for ProfileDefinition and ProfileCatalogSource
		// but then thought it might be a bit annoying when we add bits in the future?
		// To me, checking that a schema is generated is enough IDK
	})

	When("object does not exist", func() {
		It("fails", func() {
			session, err := runCmd("nothing", "")
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 20).Should(gbytes.Say("Couldn't find ref nothing in definitions"))
		})
	})

	When("the output location is a directory which does not exist", func() {
		It("fails", func() {
			destFile := filepath.Join(tmpDir, "notCreatable", "schema.json")
			session, err := runCmd(object, destFile)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 20).Should(gbytes.Say("no such file or directory"))

			Expect(destFile).ToNot(BeAnExistingFile())
		})
	})

	When("not enough args are given", func() {
		It("fails with help message", func() {
			cliArgs := []string{}
			cliCmd := exec.Command(cliBin, cliArgs...)
			session, err := gexec.Start(cliCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("Usage: schema <object-name> <outfile>"))
		})
	})
})

func runCmd(object, outFile string) (*gexec.Session, error) {
	cliArgs := []string{object, outFile}
	cliCmd := exec.Command(cliBin, cliArgs...)
	return gexec.Start(cliCmd, GinkgoWriter, GinkgoWriter)
}
