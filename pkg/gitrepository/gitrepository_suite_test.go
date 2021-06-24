package gitrepository_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGitrepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitrepository Suite")
}
