package acceptance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/protos"
)

const (
	profileAPIVersion = "profiles.weave.works/v1alpha1"
)

var _ = Describe("Acceptance", func() {
	Context("ProfileCatalog", func() {
		var (
			pCatalog                       profilesv1.ProfileCatalogSource
			expectedNginx1, expectedNginx2 profilesv1.ProfileCatalogEntry
			sourceName, profileName        string
		)

		BeforeEach(func() {
			sourceName, profileName = "catalog", "nginx-1"
			pCatalog = profilesv1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: profileAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      sourceName,
					Namespace: "default",
				},
				Spec: profilesv1.ProfileCatalogSourceSpec{
					Profiles: []profilesv1.ProfileCatalogEntry{
						{
							ProfileDescription: profilesv1.ProfileDescription{
								Name:          profileName,
								Description:   "nginx 1",
								Maintainer:    "my aunt ethel",
								Prerequisites: []string{"at least 20 years of kubernetes experience"},
							},
							URL: "foo.com/bar",
							Tag: "0.0.1",
						},
						{
							ProfileDescription: profilesv1.ProfileDescription{
								Name:          profileName,
								Description:   "nginx 1 with super cool updates",
								Maintainer:    "my latest version of aunt ethel",
								Prerequisites: []string{"at least 20 years of kubernetes experience"},
							},
							Tag: "0.0.2",
							URL: "foo.com/bar",
						},
						{
							ProfileDescription: profilesv1.ProfileDescription{
								Name:        "nginx-2",
								Description: "nginx 1",
							},
						},
						{
							ProfileDescription: profilesv1.ProfileDescription{
								Name:        "something-else",
								Description: "something else",
							},
						},
					},
				},
			}
			Expect(kClient.Create(context.Background(), &pCatalog)).To(Succeed())

			expectedNginx1 = profilesv1.ProfileCatalogEntry{
				ProfileDescription: profilesv1.ProfileDescription{
					Name:          profileName,
					Description:   "nginx 1",
					Maintainer:    "my aunt ethel",
					Prerequisites: []string{"at least 20 years of kubernetes experience"},
				},
				CatalogSource: sourceName,
				Tag:           "0.0.1",
				URL:           "foo.com/bar",
			}

			expectedNginx2 = profilesv1.ProfileCatalogEntry{
				ProfileDescription: profilesv1.ProfileDescription{
					Name:          profileName,
					Description:   "nginx 1 with super cool updates",
					Maintainer:    "my latest version of aunt ethel",
					Prerequisites: []string{"at least 20 years of kubernetes experience"},
				},
				CatalogSource: sourceName,
				Tag:           "0.0.2",
				URL:           "foo.com/bar",
			}
		})

		AfterEach(func() {
			_ = kClient.Delete(context.Background(), &pCatalog)
		})

		Context("search", func() {
			It("returns the matching catalogs", func() {
				Eventually(func() []profilesv1.ProfileCatalogEntry {
					req, err := http.NewRequest("GET", "http://localhost:8000/v1/profiles", nil)
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
					var descriptions protos.GRPCProfileCatalogEntryList
					_ = json.NewDecoder(resp.Body).Decode(&descriptions)
					return descriptions.Items
				}).Should(ConsistOf(
					expectedNginx1,
					expectedNginx2,
					profilesv1.ProfileCatalogEntry{
						ProfileDescription: profilesv1.ProfileDescription{
							Name:          "nginx-2",
							Description:   "nginx 1",
							Prerequisites: []string{},
						},
						CatalogSource: sourceName,
					},
				))
			})
		})

		Context("creating a catalog", func() {
			var catalog profilesv1.ProfileCatalogSource
			AfterEach(func() {
				_ = kClient.Delete(context.Background(), &catalog)
			})

			Context("when a valid version is provided", func() {
				It("create successfully", func() {
					catalog = profilesv1.ProfileCatalogSource{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ProfileCatalogSource",
							APIVersion: profileAPIVersion,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "invalid",
							Namespace: "default",
						},
						Spec: profilesv1.ProfileCatalogSourceSpec{
							Profiles: []profilesv1.ProfileCatalogEntry{
								{
									ProfileDescription: profilesv1.ProfileDescription{
										Name:          profileName,
										Description:   "nginx 1",
										Maintainer:    "my aunt ethel",
										Prerequisites: []string{"at least 20 years of kubernetes experience"},
									},
									Tag: "0.2.1",
									URL: "foo.com/bar",
								},
								{
									ProfileDescription: profilesv1.ProfileDescription{
										Name:          profileName,
										Description:   "nginx 1",
										Maintainer:    "my aunt ethel",
										Prerequisites: []string{"at least 20 years of kubernetes experience"},
									},
									Tag: "v0.2.1",
									URL: "foo.com/bar",
								},
								{
									ProfileDescription: profilesv1.ProfileDescription{
										Name:          profileName,
										Description:   "nginx 1",
										Maintainer:    "my aunt ethel",
										Prerequisites: []string{"at least 20 years of kubernetes experience"},
									},
									Tag: "foo-bar/0.2.1",
									URL: "foo.com/bar",
								},
								{
									ProfileDescription: profilesv1.ProfileDescription{
										Name:          profileName,
										Description:   "nginx 1",
										Maintainer:    "my aunt ethel",
										Prerequisites: []string{"at least 20 years of kubernetes experience"},
									},
									Tag: "foo-bar/v0.2.1",
									URL: "foo.com/bar",
								},
								{
									ProfileDescription: profilesv1.ProfileDescription{
										Name:          profileName,
										Description:   "nginx 1",
										Maintainer:    "my aunt ethel",
										Prerequisites: []string{"at least 20 years of kubernetes experience"},
									},
									Tag: "0.2.1-build.1",
									URL: "foo.com/bar",
								},
							},
						},
					}
					Expect(kClient.Create(context.Background(), &catalog)).To(Succeed())
				})
			})
			Context("when a invalid version is provided", func() {
				It("rejects it", func() {
					catalog = profilesv1.ProfileCatalogSource{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ProfileCatalogSource",
							APIVersion: profileAPIVersion,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "invalid",
							Namespace: "default",
						},
						Spec: profilesv1.ProfileCatalogSourceSpec{
							Profiles: []profilesv1.ProfileCatalogEntry{
								{
									ProfileDescription: profilesv1.ProfileDescription{
										Name:          profileName,
										Description:   "nginx 1",
										Maintainer:    "my aunt ethel",
										Prerequisites: []string{"at least 20 years of kubernetes experience"},
									},
									Tag: "0.not.1",
									URL: "foo.com/bar",
								},
							},
						},
					}
					Expect(kClient.Create(context.Background(), &catalog)).To(MatchError(ContainSubstring("spec.profiles.tag in body should match")))
				})
			})

			Context("when a creating from a repo URL", func() {
				var namespace string

				BeforeEach(func() {
					namespace = uuid.New().String()
					nsp := v1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: namespace,
						},
					}
					Expect(kClient.Create(context.Background(), &nsp)).To(Succeed())
				})

				AfterEach(func() {
					nsp := v1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: namespace,
						},
					}
					_ = kClient.Delete(context.Background(), &nsp)
				})

				Context("public repo", func() {
					It("works", func() {
						catalog = profilesv1.ProfileCatalogSource{
							TypeMeta: metav1.TypeMeta{
								Kind:       "ProfileCatalogSource",
								APIVersion: profileAPIVersion,
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "repo",
								Namespace: namespace,
							},
							Spec: profilesv1.ProfileCatalogSourceSpec{
								Repos: []profilesv1.Repository{
									{
										URL: "https://github.com/weaveworks/profiles-examples",
									},
								},
							},
						}

						Expect(kClient.Create(context.Background(), &catalog)).To(Succeed())
						Eventually(func() profilesv1.ProfileCatalogEntry {
							description, _ := getProfile("weaveworks-nginx", "repo", "v0.1.1")
							return description
						}, "60s", "5s").Should(Equal(profilesv1.ProfileCatalogEntry{
							ProfileDescription: profilesv1.ProfileDescription{
								Name:          "weaveworks-nginx",
								Description:   "Profile for deploying nginx",
								Maintainer:    "weaveworks",
								Prerequisites: []string{"kubernetes 1.19"},
							},
							Tag:           "weaveworks-nginx/v0.1.1",
							URL:           "https://github.com/weaveworks/profiles-examples",
							CatalogSource: "repo",
						}))
					})
				})

				Context("private repo", func() {
					var (
						secretName = "ssh-secret"
						secret     corev1.Secret
					)

					BeforeEach(func() {
						if os.Getenv("TEST_SSH_PRIVATE_KEY") == "" {
							Skip("SKIP, this test needs TEST_SSH_PRIVATE_KEY to work. You really should be running this test!")
						}
						cmd := exec.Command("ssh-keyscan", "github.com")
						knownHosts, err := cmd.CombinedOutput()
						Expect(err).NotTo(HaveOccurred())

						secret = corev1.Secret{
							TypeMeta: metav1.TypeMeta{
								Kind:       "Secret",
								APIVersion: "v1",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      secretName,
								Namespace: namespace,
							},
							Data: map[string][]byte{
								"identity":    []byte(os.Getenv("TEST_SSH_PRIVATE_KEY")),
								"known_hosts": knownHosts,
							},
						}

						Expect(kClient.Create(context.Background(), &secret)).To(Succeed())
					})

					It("works", func() {
						catalog = profilesv1.ProfileCatalogSource{
							TypeMeta: metav1.TypeMeta{
								Kind:       "ProfileCatalogSource",
								APIVersion: profileAPIVersion,
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "repo",
								Namespace: namespace,
							},
							Spec: profilesv1.ProfileCatalogSourceSpec{
								Repos: []profilesv1.Repository{
									{
										URL: "ssh://git@github.com/weaveworks/profiles-examples-private",
										SecretRef: &meta.LocalObjectReference{
											Name: secretName,
										},
									},
								},
							},
						}

						Expect(kClient.Create(context.Background(), &catalog)).To(Succeed())
						Eventually(func() profilesv1.ProfileCatalogEntry {
							description, _ := getProfile("weaveworks-nginx", "repo", "v0.2.0")
							return description
						}, "60s", "5s").Should(Equal(profilesv1.ProfileCatalogEntry{
							ProfileDescription: profilesv1.ProfileDescription{
								Name:          "weaveworks-nginx",
								Description:   "Profile for deploying nginx",
								Maintainer:    "weaveworks",
								Prerequisites: []string{"kubernetes 1.19"},
							},
							Tag:           "weaveworks-nginx/v0.2.0",
							URL:           "ssh://git@github.com/weaveworks/profiles-examples-private",
							CatalogSource: "repo",
						}))
					})
				})
			})
		})

		Context("get", func() {
			It("returns details of the requested catalog entry", func() {
				Eventually(func() profilesv1.ProfileCatalogEntry {
					description, _ := getProfile(profileName, sourceName, "")
					return description
				}, "10s").Should(Equal(expectedNginx1))
			})

			When("version is set to latest", func() {
				It("returns details of the requested catalog entry with the latest version", func() {
					Eventually(func() profilesv1.ProfileCatalogEntry {
						description, _ := getProfile(profileName, sourceName, "latest")
						return description
					}, "10s").Should(Equal(expectedNginx2))
				})
			})

			When("a request is made to list all available updates", func() {
				It("returns a list of available profiles with greater versions", func() {
					Eventually(func() []profilesv1.ProfileCatalogEntry {
						versions, err := getProfileUpdates(profileName, sourceName, "0.0.1")
						Expect(err).NotTo(HaveOccurred())
						return versions
					}, "10s").Should(ContainElement(expectedNginx2))
				})
			})
		})

		Context("update", func() {
			It("updates a ProfileCatalogSource with new profiles", func() {
				pCatalog.Spec.Profiles = append(pCatalog.Spec.Profiles, profilesv1.ProfileCatalogEntry{
					ProfileDescription: profilesv1.ProfileDescription{
						Name:        "new-profile",
						Description: "I am new here",
					},
				})
				Expect(kClient.Update(context.Background(), &pCatalog)).To(Succeed())
				Eventually(func() profilesv1.ProfileCatalogEntry {
					description, err := getProfile("new-profile", sourceName, "")
					Expect(err).NotTo(HaveOccurred())
					return description
				}).Should(Equal(profilesv1.ProfileCatalogEntry{
					ProfileDescription: profilesv1.ProfileDescription{
						Name:          "new-profile",
						Description:   "I am new here",
						Prerequisites: []string{},
					},
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

func getProfile(profileName, sourceName, version string) (profilesv1.ProfileCatalogEntry, error) {
	u, err := url.Parse("http://localhost:8000/v1/profiles")
	if err != nil {
		return profilesv1.ProfileCatalogEntry{}, err
	}
	u.Path = path.Join(u.Path, sourceName, profileName, version)
	resp, err := http.Get(u.String())
	if err != nil {
		return profilesv1.ProfileCatalogEntry{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return profilesv1.ProfileCatalogEntry{}, fmt.Errorf("expected status code 200; got %d", resp.StatusCode)
	}
	var p protos.GRPCProfileCatalogEntry
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return profilesv1.ProfileCatalogEntry{}, err
	}
	return p.Item, nil
}

func getProfileUpdates(profileName, sourceName, version string) ([]profilesv1.ProfileCatalogEntry, error) {
	u := fmt.Sprintf("http://localhost:8000/v1/profiles/%s/%s/%s/available_updates", sourceName, profileName, version)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200; got %d", resp.StatusCode)
	}
	var p protos.GRPCProfileCatalogEntryList
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return p.Items, nil
}
