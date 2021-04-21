package profile

import (
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/extensions/table"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm", func() {
	Describe("HelmRepo", func() {
		var oldRes, newRes *sourcev1.HelmRepository
		BeforeEach(func() {
			oldRes = &sourcev1.HelmRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       helmv2.HelmReleaseKind,
					APIVersion: helmv2.GroupVersion.String(),
				},
				Spec: sourcev1.HelmRepositorySpec{
					URL: "example.com",
				},
			}
			newRes = oldRes.DeepCopy()
		})

		DescribeTable("helmReleaseResourceRequiresUpdate",
			func(newRes func() *sourcev1.HelmRepository, updateExpected bool) {
				Expect(helmRepoRequiresUpdate(oldRes, newRes())).To(Equal(updateExpected))
			},
			Entry("spec is unchanged should return false", func() *sourcev1.HelmRepository {
				return newRes
			}, false),
			Entry("chart change should return true", func() *sourcev1.HelmRepository {
				newRes.Spec.URL = "new"
				return newRes
			}, true),
		)

	})
	Describe("HelmRelease", func() {
		var oldRes, newRes *helmv2.HelmRelease
		BeforeEach(func() {
			oldRes = &helmv2.HelmRelease{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       helmv2.HelmReleaseKind,
					APIVersion: helmv2.GroupVersion.String(),
				},
				Spec: helmv2.HelmReleaseSpec{
					Chart: helmv2.HelmChartTemplate{
						Spec: helmv2.HelmChartTemplateSpec{
							Chart: "foo",
							SourceRef: helmv2.CrossNamespaceObjectReference{
								Kind:      sourcev1.HelmRepositoryKind,
								Name:      "repo",
								Namespace: "default",
							},
							Version: "1.0",
						},
					},
					Values: &apiextensionsv1.JSON{Raw: []byte("foo")},
					ValuesFrom: []helmv2.ValuesReference{
						{
							Name:     "nginx-values",
							Kind:     "Secret",
							Optional: true,
						},
					},
				},
			}
			newRes = oldRes.DeepCopy()
		})

		DescribeTable("helmReleaseResourceRequiresUpdate",
			func(newRes func() *helmv2.HelmRelease, updateExpected bool) {
				Expect(helmReleaseRequiresUpdate(oldRes, newRes())).To(Equal(updateExpected))
			},
			Entry("spec is unchanged should return false", func() *helmv2.HelmRelease {
				return newRes
			}, false),
			Entry("chart change should return true", func() *helmv2.HelmRelease {
				newRes.Spec.Chart.Spec.Chart = "new"
				return newRes
			}, true),
			Entry("sourceRef change should return true", func() *helmv2.HelmRelease {
				newRes.Spec.Chart.Spec.SourceRef.Name = "new"
				return newRes
			}, true),
			Entry("version change should return true", func() *helmv2.HelmRelease {
				newRes.Spec.Chart.Spec.Version = "new"
				return newRes
			}, true),
			Entry("values change should return true", func() *helmv2.HelmRelease {
				newRes.Spec.Values = &apiextensionsv1.JSON{Raw: []byte("new")}
				return newRes
			}, true),
			Entry("valuesFrom change should return true", func() *helmv2.HelmRelease {
				newRes.Spec.ValuesFrom = []helmv2.ValuesReference{
					{
						Name:     "new",
						Kind:     "Secret",
						Optional: true,
					},
				}
				return newRes
			}, true),
		)
	})
})
