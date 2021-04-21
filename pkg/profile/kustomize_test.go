package profile

import (
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/extensions/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kustomization", func() {
	var oldRes, newRes *kustomizev1.Kustomization
	BeforeEach(func() {
		oldRes = &kustomizev1.Kustomization{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       kustomizev1.KustomizationKind,
				APIVersion: kustomizev1.GroupVersion.String(),
			},
			Spec: kustomizev1.KustomizationSpec{
				Path:            "path/1",
				Interval:        metav1.Duration{Duration: time.Minute * 5},
				Prune:           true,
				TargetNamespace: "default",
				SourceRef: kustomizev1.CrossNamespaceSourceReference{
					Kind:      sourcev1.GitRepositoryKind,
					Name:      "foo",
					Namespace: "default",
				},
			},
		}
		newRes = oldRes.DeepCopy()
	})

	DescribeTable("kustomizationResourceRequiresUpdate",
		func(newRes func() *kustomizev1.Kustomization, updateExpected bool) {
			Expect(kustomizeRequiresUpdate(oldRes, newRes())).To(Equal(updateExpected))
		},
		Entry("spec is unchanged should return false", func() *kustomizev1.Kustomization {
			return newRes
		}, false),
		Entry("path change should return true", func() *kustomizev1.Kustomization {
			newRes.Spec.Path = "new/path"
			return newRes
		}, true),
		Entry("interval change should return true", func() *kustomizev1.Kustomization {
			newRes.Spec.Interval = metav1.Duration{Duration: time.Minute * 4}
			return newRes
		}, true),
		Entry("prune change should return true", func() *kustomizev1.Kustomization {
			newRes.Spec.Prune = false
			return newRes
		}, true),
		Entry("targetNamespace change should return true", func() *kustomizev1.Kustomization {
			newRes.Spec.TargetNamespace = "new"
			return newRes
		}, true),
		Entry("sourceRef change should return true", func() *kustomizev1.Kustomization {
			newRes.Spec.SourceRef.Name = "new"
			return newRes
		}, true),
	)
})
