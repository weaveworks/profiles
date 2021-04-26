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

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`[{"name":"nginx-1","description":"nginx 1","catalog":"foo"}]`))
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

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`[]`))
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

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`{"name":"nginx-1","description":"nginx 1","catalog":"catalog"}`))
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

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
