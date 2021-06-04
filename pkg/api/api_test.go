package api_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/api"
	catfakes "github.com/weaveworks/profiles/pkg/api/fakes"
)

var _ = Describe("Api", func() {
	var (
		catalogAPI  api.API
		fakeCatalog *catfakes.FakeCatalog
	)

	BeforeEach(func() {
		fakeCatalog = new(catfakes.FakeCatalog)
		catalogAPI = api.New(fakeCatalog)
	})

	Context("/profiles", func() {
		When("a matching profile exists", func() {
			BeforeEach(func() {
				fakeCatalog.SearchReturns([]profilesv1.ProfileDescription{
					{
						Name:          "nginx-1",
						Description:   "nginx 1",
						CatalogSource: "foo",
					},
				})
			})

			It("returns the matching profiles from the catalog", func() {
				req, err := http.NewRequest("GET", "/profiles", nil)
				Expect(err).NotTo(HaveOccurred())
				u, err := url.Parse("http://example.com")
				Expect(err).NotTo(HaveOccurred())
				q := u.Query()
				q.Add("name", "nginx")
				req.URL.RawQuery = q.Encode()
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfilesHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`[{"name":"nginx-1","description":"nginx 1","catalog":"foo"}]`))
				Expect(fakeCatalog.SearchCallCount()).To(Equal(1))
				actualProfileName := fakeCatalog.SearchArgsForCall(0)
				Expect(actualProfileName).To(Equal("nginx"))
			})
		})

		When("a no matching profiles are found", func() {
			BeforeEach(func() {
				fakeCatalog.SearchReturns([]profilesv1.ProfileDescription{})
			})

			It("returns an empty array but does not 404", func() {
				req, err := http.NewRequest("GET", "/profiles", nil)
				Expect(err).NotTo(HaveOccurred())
				u, err := url.Parse("http://example.com")
				Expect(err).NotTo(HaveOccurred())
				q := u.Query()
				q.Add("name", "nginx")
				req.URL.RawQuery = q.Encode()
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfilesHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`[]`))
				Expect(fakeCatalog.SearchCallCount()).To(Equal(1))
				actualProfileName := fakeCatalog.SearchArgsForCall(0)
				Expect(actualProfileName).To(Equal("nginx"))
			})
		})

		When("no name query is provided", func() {
			BeforeEach(func() {
				fakeCatalog.SearchReturns([]profilesv1.ProfileDescription{})
			})

			It("returns a 400", func() {
				req, err := http.NewRequest("GET", "/profiles", nil)
				Expect(err).NotTo(HaveOccurred())
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfilesHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
				Expect(fakeCatalog.SearchCallCount()).To(Equal(0))
			})
		})
	})

	Context("/profiles/catalog/profile-name", func() {
		var (
			sourceName, profileName string
		)

		BeforeEach(func() {
			sourceName, profileName = "catalog", "nginx-1"
		})

		When("the requested profile exists", func() {
			BeforeEach(func() {
				fakeCatalog.GetReturns(&profilesv1.ProfileDescription{
					Name:          "nginx-1",
					Description:   "nginx 1",
					CatalogSource: "catalog",
				})
			})

			It("returns the profile summary from the catalog", func() {
				req, err := http.NewRequest("GET", "/profiles", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "profile": profileName})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`{"name":"nginx-1","description":"nginx 1","catalog":"catalog"}`))
				Expect(fakeCatalog.GetCallCount()).To(Equal(1))
				actualSourceName, actualProfileName := fakeCatalog.GetArgsForCall(0)
				Expect(actualSourceName).To(Equal(sourceName))
				Expect(actualProfileName).To(Equal(profileName))
			})
		})

		When("the requested profile does not exist", func() {
			BeforeEach(func() {
				fakeCatalog.GetReturns(nil)
			})

			It("returns a 404", func() {
				req, err := http.NewRequest("GET", "/profiles", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "profile": profileName})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
				Expect(fakeCatalog.GetCallCount()).To(Equal(1))
				actualSourceName, actualProfileName := fakeCatalog.GetArgsForCall(0)
				Expect(actualSourceName).To(Equal(sourceName))
				Expect(actualProfileName).To(Equal(profileName))
			})
		})

		When("a query param is missing", func() {
			BeforeEach(func() {
				fakeCatalog.GetReturns(nil)
			})
			Context("profileName", func() {
				It("returns a 400", func() {
					req, err := http.NewRequest("GET", "/profiles", nil)
					req = mux.SetURLVars(req, map[string]string{"catalog": sourceName})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.GetCallCount()).To(Equal(0))
				})
			})

			Context("profileName", func() {
				It("returns a 400", func() {
					req, err := http.NewRequest("GET", "/profiles", nil)
					req = mux.SetURLVars(req, map[string]string{"profileName": sourceName})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.GetCallCount()).To(Equal(0))
				})
			})
		})
	})

	Context("/profiles/catalog/profile-name/version", func() {
		var (
			sourceName, profileName, version string
		)

		BeforeEach(func() {
			sourceName, profileName, version = "catalog", "nginx-1", "v0.1.0"
		})

		When("the requested profile exists", func() {
			BeforeEach(func() {
				fakeCatalog.GetWithVersionReturns(&profilesv1.ProfileDescription{
					Name:          "nginx-1",
					Description:   "nginx 1",
					CatalogSource: "catalog",
					Version:       "v0.1.0",
				})
			})

			It("returns the profile summary from the catalog", func() {
				req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.1.0", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "profile": profileName, "version": version})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileWithVersionHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`{"name":"nginx-1","description":"nginx 1","version":"v0.1.0","catalog":"catalog"}`))

				Expect(fakeCatalog.GetWithVersionCallCount()).To(Equal(1))
				actualSourceName, actualProfileName, actualCatalogVersion := fakeCatalog.GetWithVersionArgsForCall(0)
				Expect(actualSourceName).To(Equal(sourceName))
				Expect(actualProfileName).To(Equal(profileName))
				Expect(actualCatalogVersion).To(Equal(version))
			})
		})

		When("the requested profile does not exist", func() {
			BeforeEach(func() {
				fakeCatalog.GetWithVersionReturns(nil)
			})

			It("returns a 404", func() {
				req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "profile": profileName, "version": version})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileWithVersionHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
				Expect(fakeCatalog.GetWithVersionCallCount()).To(Equal(1))
				actualSourceName, actualProfileName, actualCatalogVersion := fakeCatalog.GetWithVersionArgsForCall(0)
				Expect(actualSourceName).To(Equal(sourceName))
				Expect(actualProfileName).To(Equal(profileName))
				Expect(actualCatalogVersion).To(Equal(version))
			})
		})

		When("a querry param is missing", func() {
			BeforeEach(func() {
				fakeCatalog.GetWithVersionReturns(nil)
			})

			Context("catalog param", func() {
				It("returns a 404", func() {
					req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0", nil)
					req = mux.SetURLVars(req, map[string]string{"profile": profileName, "version": version})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileWithVersionHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.GetWithVersionCallCount()).To(Equal(0))
				})
			})

			Context("profile param", func() {
				It("returns a 404", func() {
					req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0", nil)
					req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "version": version})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileWithVersionHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.GetWithVersionCallCount()).To(Equal(0))
				})
			})

			Context("version param", func() {
				It("returns a 404", func() {
					req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0", nil)
					req = mux.SetURLVars(req, map[string]string{"profile": profileName, "catalog": sourceName})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileWithVersionHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.GetWithVersionCallCount()).To(Equal(0))
				})
			})
		})
	})

	Context("/profiles/catalog/profile-name/version/available_updates", func() {
		var (
			sourceName, profileName, version string
		)

		BeforeEach(func() {
			sourceName, profileName, version = "catalog", "nginx-1", "v0.1.0"
		})

		When("the requested profile has newer versions", func() {
			BeforeEach(func() {
				fakeCatalog.ProfilesGreaterThanVersionReturns([]profilesv1.ProfileDescription{
					{
						Name:          "nginx-1",
						Description:   "nginx 1",
						CatalogSource: "catalog",
						Version:       "v0.1.1",
					},
				})
			})

			It("returns the profiles with newer versions", func() {
				req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.1.0/available_updates", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "profile": profileName, "version": version})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileGreaterThanVersionHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`[{"name":"nginx-1","description":"nginx 1","version":"v0.1.1","catalog":"catalog"}]`))

				Expect(fakeCatalog.ProfilesGreaterThanVersionCallCount()).To(Equal(1))
				actualSourceName, actualProfileName, actualCatalogVersion := fakeCatalog.ProfilesGreaterThanVersionArgsForCall(0)
				Expect(actualSourceName).To(Equal(sourceName))
				Expect(actualProfileName).To(Equal(profileName))
				Expect(actualCatalogVersion).To(Equal(version))
			})
		})

		When("the requested profile does not exist", func() {
			BeforeEach(func() {
				fakeCatalog.ProfilesGreaterThanVersionReturns([]profilesv1.ProfileDescription{})
			})

			It("returns a 404", func() {
				req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0/available_updates", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "profile": profileName, "version": version})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileGreaterThanVersionHandler)

				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
				Expect(fakeCatalog.ProfilesGreaterThanVersionCallCount()).To(Equal(1))
				actualSourceName, actualProfileName, actualCatalogVersion := fakeCatalog.ProfilesGreaterThanVersionArgsForCall(0)
				Expect(actualSourceName).To(Equal(sourceName))
				Expect(actualProfileName).To(Equal(profileName))
				Expect(actualCatalogVersion).To(Equal(version))
			})
		})

		When("a querry param is missing", func() {
			BeforeEach(func() {
				fakeCatalog.GetWithVersionReturns(nil)
			})

			Context("catalog param", func() {
				It("returns a 404", func() {
					req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0/available_updates", nil)
					req = mux.SetURLVars(req, map[string]string{"profile": profileName, "version": version})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileGreaterThanVersionHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.ProfilesGreaterThanVersionCallCount()).To(Equal(0))
				})
			})

			Context("profile param", func() {
				It("returns a 404", func() {
					req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0/available_updates", nil)
					req = mux.SetURLVars(req, map[string]string{"catalog": sourceName, "version": version})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileGreaterThanVersionHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.ProfilesGreaterThanVersionCallCount()).To(Equal(0))
				})
			})

			Context("version param", func() {
				It("returns a 404", func() {
					req, err := http.NewRequest("GET", "/profile/catalog/nginx-1/v0.3.0/available_updates", nil)
					req = mux.SetURLVars(req, map[string]string{"profile": profileName, "catalog": sourceName})
					Expect(err).NotTo(HaveOccurred())

					rr := httptest.NewRecorder()
					handler := http.HandlerFunc(catalogAPI.ProfileGreaterThanVersionHandler)

					handler.ServeHTTP(rr, req)

					Expect(rr.Code).To(Equal(http.StatusBadRequest))
					Expect(fakeCatalog.ProfilesGreaterThanVersionCallCount()).To(Equal(0))
				})
			})
		})
	})
})
