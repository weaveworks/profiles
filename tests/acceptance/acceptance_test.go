package acceptance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

const (
	profileSubscriptionAPIVersion = "profilesubscriptions.weave.works/v1alpha1"
)

var _ = Describe("Acceptance", func() {
	Context("ProfileCatalog", func() {
		var (
			pCatalog                       profilesv1.ProfileCatalogSource
			expectedNginx1, expectedNginx2 profilesv1.ProfileDescription
			sourceName, profileName        string
		)

		BeforeEach(func() {
			sourceName, profileName = "catalog", "nginx-1"
			pCatalog = profilesv1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: profileSubscriptionAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      sourceName,
					Namespace: "default",
				},
				Spec: profilesv1.ProfileCatalogSourceSpec{
					Profiles: []profilesv1.ProfileDescription{
						{
							Name:          profileName,
							Description:   "nginx 1",
							Version:       "0.0.1",
							URL:           "foo.com/bar",
							Maintainer:    "my aunt ethel",
							Prerequisites: []string{"at least 20 years of kubernetes experience"},
						},
						{
							Name:          profileName,
							Description:   "nginx 1 with super cool updates",
							Version:       "0.0.2",
							URL:           "foo.com/bar",
							Maintainer:    "my latest version of aunt ethel",
							Prerequisites: []string{"at least 20 years of kubernetes experience"},
						},
						{
							Name:        "nginx-2",
							Description: "nginx 1",
						},
						{
							Name:        "something-else",
							Description: "something else",
						},
					},
				},
			}
			Expect(kClient.Create(context.Background(), &pCatalog)).To(Succeed())

			expectedNginx1 = profilesv1.ProfileDescription{
				Name:          profileName,
				Description:   "nginx 1",
				CatalogSource: sourceName,
				Version:       "0.0.1",
				URL:           "foo.com/bar",
				Maintainer:    "my aunt ethel",
				Prerequisites: []string{"at least 20 years of kubernetes experience"},
			}

			expectedNginx2 = profilesv1.ProfileDescription{
				Name:          profileName,
				Description:   "nginx 1 with super cool updates",
				CatalogSource: sourceName,
				Version:       "0.0.2",
				URL:           "foo.com/bar",
				Maintainer:    "my latest version of aunt ethel",
				Prerequisites: []string{"at least 20 years of kubernetes experience"},
			}
		})

		AfterEach(func() {
			_ = kClient.Delete(context.Background(), &pCatalog)
		})

		Context("search", func() {
			It("returns the matching catalogs", func() {
				Eventually(func() []profilesv1.ProfileDescription {
					req, err := http.NewRequest("GET", "http://localhost:8000/profiles", nil)
					Expect(err).NotTo(HaveOccurred())
					u, err := url.Parse("http://localhost:8000")
					Expect(err).NotTo(HaveOccurred())
					q := u.Query()
					q.Add("name", "nginx")
					req.URL.RawQuery = q.Encode()
					Expect(err).NotTo(HaveOccurred())
					resp, err := http.DefaultClient.Do(req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					descriptions := []profilesv1.ProfileDescription{}
					_ = json.NewDecoder(resp.Body).Decode(&descriptions)
					return descriptions
				}).Should(ConsistOf(
					expectedNginx1,
					expectedNginx2,
					profilesv1.ProfileDescription{
						Name:          "nginx-2",
						Description:   "nginx 1",
						CatalogSource: sourceName,
					},
				))
			})
		})

		Context("get", func() {
			It("returns details of the requested catalog entry", func() {
				Eventually(func() profilesv1.ProfileDescription {
					description, _ := getProfile(profileName, sourceName, "")
					return description
				}, "10s").Should(Equal(expectedNginx1))
			})

			When("version is set to latest", func() {
				It("returns details of the requested catalog entry with the latest version", func() {
					Eventually(func() profilesv1.ProfileDescription {
						description, _ := getProfile(profileName, sourceName, "latest")
						return description
					}, "10s").Should(Equal(expectedNginx2))
				})
			})
		})

		Context("update", func() {
			It("updates a ProfileCatalogSource with new profiles", func() {
				pCatalog.Spec.Profiles = append(pCatalog.Spec.Profiles, profilesv1.ProfileDescription{
					Name:        "new-profile",
					Description: "I am new here",
				})
				Expect(kClient.Update(context.Background(), &pCatalog)).To(Succeed())
				Eventually(func() profilesv1.ProfileDescription {
					description, err := getProfile("new-profile", sourceName, "")
					Expect(err).NotTo(HaveOccurred())
					return description
				}).Should(Equal(profilesv1.ProfileDescription{
					Name:          "new-profile",
					Description:   "I am new here",
					CatalogSource: sourceName,
				}))
			})
		})

		Context("delete", func() {
			It("clears the in-memory cache when a ProfileCatalogSource is deleted", func() {
				description, err := getProfile(profileName, sourceName, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(description).To(Equal(expectedNginx1))

				Expect(kClient.Delete(context.Background(), &pCatalog)).To(Succeed())
				Eventually(func() error {
					_, err := getProfile(profileName, sourceName, "")
					return err
				}, "5s").Should(MatchError(ContainSubstring("got 404")))
			})
		})
	})
})

func getProfile(profileName, sourceName, version string) (profilesv1.ProfileDescription, error) {
	url := fmt.Sprintf("http://localhost:8000/profiles/%s/%s", sourceName, profileName)
	if version != "" {
		url = fmt.Sprintf("%s/%s", url, version)
	}
	resp, err := http.Get(url)
	if err != nil {
		return profilesv1.ProfileDescription{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return profilesv1.ProfileDescription{}, fmt.Errorf("expected status code 200; got %d", resp.StatusCode)
	}
	var p profilesv1.ProfileDescription
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return profilesv1.ProfileDescription{}, err
	}
	return p, nil
}
