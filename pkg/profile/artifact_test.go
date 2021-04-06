package profile_test

import (
	"context"
	"fmt"
	"reflect"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/profile"
)

const (
	subscriptionName     = "mySub"
	namespace            = "default"
	branch               = "main"
	profileName          = "profileName"
	chartName            = "artifactName"
	chartPath            = "artifactPath"
	profileSubKind       = "ProfileSubscription"
	profileSubAPIVersion = "weave.works/v1alpha1"
)

var (
	profileTypeMeta = metav1.TypeMeta{
		Kind:       profileSubKind,
		APIVersion: profileSubAPIVersion,
	}
)

var _ = Describe("Profile", func() {
	var (
		ctx        = context.Background()
		p          *profile.Profile
		fakeClient client.Client
		scheme     *runtime.Scheme
		pSub       v1alpha1.ProfileSubscription
		pDef       v1alpha1.ProfileDefinition
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		pSub = v1alpha1.ProfileSubscription{
			TypeMeta: profileTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      subscriptionName,
				Namespace: namespace,
			},
			Spec: v1alpha1.ProfileSubscriptionSpec{
				ProfileURL: "https://github.com/org/repo-name",
				Branch:     branch,
				Values: &apiextensionsv1.JSON{
					Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
				},
				ValuesFrom: []helmv2.ValuesReference{
					{
						Name:     "nginx-values",
						Kind:     "Secret",
						Optional: true,
					},
				},
			},
		}

		pDef = v1alpha1.ProfileDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: profileName,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "Profile",
				APIVersion: "profiles.fluxcd.io/v1alpha1",
			},
			Spec: v1alpha1.ProfileDefinitionSpec{
				Description: "foo",
				Artifacts: []v1alpha1.Artifact{
					{
						Name: chartName,
						Path: chartPath,
						Kind: "HelmChart",
					},
				},
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
	})

	var _ = Describe("MakeArtifacts", func() {
		It("creates a slice of runtime.Object", func() {
			Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
			Expect(helmv2.AddToScheme(scheme)).To(Succeed())
			Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())

			p = profile.New(pDef, pSub, fakeClient, logr.Discard())
			o, err := p.MakeArtifacts()
			Expect(err).NotTo(HaveOccurred())

			Expect(o).To(HaveLen(2))
			Expect(o[0]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "GitRepository", APIVersion: "source.toolkit.fluxcd.io/v1beta1"}))
			Expect(o[1]).To(HaveTypeMeta(metav1.TypeMeta{Kind: "HelmRelease", APIVersion: "helm.toolkit.fluxcd.io/v2beta1"}))
		})
	})

	var _ = Describe("CreateArtifacts", func() {
		When("the profile consists of a HelmChart", func() {
			It("creates the helm and gitrepo resources with the correct owner", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())
				Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())

				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).NotTo(HaveOccurred())

				gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "repo-name", branch)
				gitRepo := sourcev1.GitRepository{}
				err = fakeClient.Get(ctx, client.ObjectKey{Name: gitRefName, Namespace: namespace}, &gitRepo)
				Expect(err).NotTo(HaveOccurred())
				Expect(gitRepo.Spec.URL).To(Equal("https://github.com/org/repo-name"))
				Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))
				Expect(gitRepo.OwnerReferences).To(HaveLen(1))
				Expect(gitRepo.OwnerReferences[0].Name).To(Equal(subscriptionName))
				Expect(gitRepo.OwnerReferences[0].Kind).To(Equal(profileSubKind))
				Expect(gitRepo.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
				Expect(*gitRepo.OwnerReferences[0].Controller).To(BeTrue())

				helmReleaseName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName)
				helmRelease := helmv2.HelmRelease{}
				err = fakeClient.Get(ctx, client.ObjectKey{Name: helmReleaseName, Namespace: namespace}, &helmRelease)
				Expect(err).NotTo(HaveOccurred())
				Expect(helmRelease.Spec.Chart.Spec.Chart).To(Equal(chartPath))
				Expect(helmRelease.Spec.Chart.Spec.SourceRef).To(Equal(
					helmv2.CrossNamespaceObjectReference{
						Kind:      "GitRepository",
						Name:      gitRefName,
						Namespace: namespace,
					},
				))
				Expect(helmRelease.GetValues()).To(Equal(map[string]interface{}{
					"replicaCount": float64(3),
					"service": map[string]interface{}{
						"port": float64(8081),
					},
				}))
				Expect(helmRelease.Spec.ValuesFrom).To(Equal([]helmv2.ValuesReference{
					{
						Name:     "nginx-values",
						Kind:     "Secret",
						Optional: true,
					},
				}))
				Expect(helmRelease.OwnerReferences).To(HaveLen(1))
				Expect(helmRelease.OwnerReferences[0].Name).To(Equal(subscriptionName))
				Expect(helmRelease.OwnerReferences[0].Kind).To(Equal(profileSubKind))
				Expect(helmRelease.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
				Expect(*helmRelease.OwnerReferences[0].Controller).To(BeTrue())
			})
		})

		When("the profile consists of yaml", func() {
			BeforeEach(func() {
				pDef.Spec.Artifacts[0].Kind = "Kustomize"
			})

			It("creates the kustomize and gitrepo resources with the correct owner", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(kustomizev1.AddToScheme(scheme)).To(Succeed())
				Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())

				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).NotTo(HaveOccurred())

				gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "repo-name", branch)
				gitRepo := sourcev1.GitRepository{}
				err = fakeClient.Get(ctx, client.ObjectKey{Name: gitRefName, Namespace: namespace}, &gitRepo)
				Expect(err).NotTo(HaveOccurred())
				Expect(gitRepo.Spec.URL).To(Equal("https://github.com/org/repo-name"))
				Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))
				Expect(gitRepo.OwnerReferences).To(HaveLen(1))
				Expect(gitRepo.OwnerReferences[0].Name).To(Equal(subscriptionName))
				Expect(gitRepo.OwnerReferences[0].Kind).To(Equal(profileSubKind))
				Expect(gitRepo.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
				Expect(*gitRepo.OwnerReferences[0].Controller).To(BeTrue())

				kustomizeName := fmt.Sprintf("%s-%s-%s", subscriptionName, profileName, chartName)
				kustomize := kustomizev1.Kustomization{}
				err = fakeClient.Get(ctx, client.ObjectKey{Name: kustomizeName, Namespace: namespace}, &kustomize)
				Expect(err).NotTo(HaveOccurred())
				Expect(kustomize.Spec.Path).To(Equal(chartPath))
				Expect(kustomize.Spec.TargetNamespace).To(Equal(namespace))
				Expect(kustomize.Spec.Prune).To(BeTrue())
				Expect(kustomize.Spec.Interval).To(Equal(metav1.Duration{Duration: time.Minute * 5}))
				Expect(kustomize.Spec.SourceRef).To(Equal(
					kustomizev1.CrossNamespaceSourceReference{
						Kind:      "GitRepository",
						Name:      gitRefName,
						Namespace: namespace,
					},
				))
				Expect(kustomize.OwnerReferences).To(HaveLen(1))
				Expect(kustomize.OwnerReferences[0].Name).To(Equal(subscriptionName))
				Expect(kustomize.OwnerReferences[0].Kind).To(Equal(profileSubKind))
				Expect(kustomize.OwnerReferences[0].APIVersion).To(Equal(profileSubAPIVersion))
				Expect(*kustomize.OwnerReferences[0].Controller).To(BeTrue())
			})
		})

		When("setting the resource owner fails", func() {
			It("errors", func() {
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to set resource ownership")))
			})
		})

		When("the GitRepository create fails", func() {
			It("errors", func() {
				// this is a bit of a hack, but by not adding this resource to the scheme
				// we can force the Create call to fail
				Expect(helmv2.AddToScheme(scheme)).To(Succeed())
				Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to create GitRepository resource")))
			})
		})

		When("the HelmRelease create fails", func() {
			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to create HelmRelease resource")))
			})
		})

		When("the Kustomization create fails", func() {
			BeforeEach(func() {
				pDef.Spec.Artifacts[0].Kind = "Kustomize"
			})

			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("failed to create Kustomization resource")))
			})
		})

		When("the Kind of artifact is unknown", func() {
			BeforeEach(func() {
				pDef.Spec.Artifacts[0].Kind = "SomeUnknownKind"
			})

			It("errors", func() {
				Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
				Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
				p = profile.New(pDef, pSub, fakeClient, logr.Discard())
				err := p.CreateArtifacts(ctx)
				Expect(err).To(MatchError(ContainSubstring("artifact kind \"SomeUnknownKind\" not recognized")))
			})
		})
	})
})

func HaveTypeMeta(expected metav1.TypeMeta) types.GomegaMatcher {
	return &typeMetaMatcher{
		expected: expected,
	}
}

type typeMetaMatcher struct {
	expected metav1.TypeMeta
}

func (m *typeMetaMatcher) Match(actual interface{}) (bool, error) {
	ro, ok := actual.(runtime.Object)
	if !ok {
		return false, fmt.Errorf("HaveTypeMeta expects a runtime.Object")
	}

	tm, err := m.typeMetaFromObject(ro)
	if err != nil {
		return false, fmt.Errorf("failed to get the type meta for object %#v: %w", ro, err)
	}

	return reflect.DeepEqual(tm, m.expected), nil
}

func (m *typeMetaMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nto have TypeMeta\n\t%#v", actual, m.expected)
}

func (m *typeMetaMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nnot to have TypeMeta\n\t%#v", actual, m.expected)
}

func (m *typeMetaMatcher) typeMetaFromObject(ro runtime.Object) (metav1.TypeMeta, error) {
	ta, err := meta.TypeAccessor(ro)
	if err != nil {
		return metav1.TypeMeta{}, fmt.Errorf("failed to get the type meta for object %#v: %w", ro, err)
	}
	return metav1.TypeMeta{APIVersion: ta.GetAPIVersion(), Kind: ta.GetKind()}, nil
}
