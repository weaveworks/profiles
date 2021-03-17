package profile_test

import (
	"context"
	"fmt"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/profile"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Artifact", func() {
	var (
		ctx              = context.Background()
		subscriptionName string
		namespace        string
		branch           string
		profileName      string
		chartName        string
		chartPath        string

		p          *profile.Profile
		fakeClient client.Client

		scheme *runtime.Scheme
	)

	BeforeEach(func() {
		subscriptionName = "mySub"
		namespace = "default"
		branch = "main"
		profileName = "profileName"
		chartName = "artifactName"
		chartPath = "artifactPath"

		pSub := v1alpha1.ProfileSubscription{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ProfileSubscription",
				APIVersion: "profilesubscriptions.weave.works/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      subscriptionName,
				Namespace: namespace,
			},
			Spec: v1alpha1.ProfileSubscriptionSpec{
				ProfileURL: "https://github.com/org/repo-name",
				Branch:     branch,
			},
		}

		pDef := v1alpha1.ProfileDefinition{
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
					},
				},
			},
		}

		scheme = runtime.NewScheme()
		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
		p = profile.New(pDef, pSub, fakeClient, logr.Discard())
	})

	It("creates the helm and gitrepo resources", func() {
		Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
		Expect(helmv2.AddToScheme(scheme)).To(Succeed())

		err := p.CreateArtifacts(ctx)
		Expect(err).NotTo(HaveOccurred())

		gitRefName := fmt.Sprintf("%s-%s-%s", subscriptionName, "repo-name", branch)
		gitRepo := sourcev1.GitRepository{}
		err = fakeClient.Get(ctx, client.ObjectKey{Name: gitRefName, Namespace: namespace}, &gitRepo)
		Expect(err).NotTo(HaveOccurred())
		Expect(gitRepo.Spec.URL).To(Equal("https://github.com/org/repo-name"))
		Expect(gitRepo.Spec.Reference.Branch).To(Equal(branch))

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
	})

	When("the GitRepository create fails", func() {
		It("errors", func() {
			// this is a bit of a hack, but by not adding this resource to the scheme
			// we can force the Create call to fail
			Expect(helmv2.AddToScheme(scheme)).To(Succeed())
			err := p.CreateArtifacts(ctx)
			Expect(err).To(MatchError(ContainSubstring("failed to create GitRepository resource")))
		})
	})

	When("the HelmRelease create fails", func() {
		It("errors", func() {
			Expect(sourcev1.AddToScheme(scheme)).To(Succeed())
			err := p.CreateArtifacts(ctx)
			Expect(err).To(MatchError(ContainSubstring("failed to create HelmRelease resource")))
		})
	})
})
