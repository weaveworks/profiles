package profile

import (
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/extensions/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git", func() {
	var oldRes, newRes *sourcev1.GitRepository
	BeforeEach(func() {
		oldRes = &sourcev1.GitRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "git",
				Namespace: "default",
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       sourcev1.GitRepositoryKind,
				APIVersion: sourcev1.GroupVersion.String(),
			},
			Spec: sourcev1.GitRepositorySpec{
				URL: "example.com",
				Reference: &sourcev1.GitRepositoryRef{
					Branch: "main",
				},
			},
		}
		newRes = oldRes.DeepCopy()
	})

	DescribeTable("gitResourceRequiresUpdate",
		func(newRes func() *sourcev1.GitRepository, updateExpected bool) {
			Expect(gitRepoRequiresUpdate(oldRes, newRes())).To(Equal(updateExpected))
		},
		Entry("spec is unchanged should return false", func() *sourcev1.GitRepository {
			return newRes
		}, false),
		Entry("url change should return true", func() *sourcev1.GitRepository {
			newRes.Spec.URL = "example.com/something"
			return newRes
		}, true),
		Entry("branch change should return true", func() *sourcev1.GitRepository {
			newRes.Spec.Reference.Branch = "foo"
			return newRes
		}, true),
	)
})
