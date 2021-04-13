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
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Api", func() {
	var (
		catalogAPI     api.API
		profileCatalog catalog.Catalog
	)
	Context("/profiles", func() {
		BeforeEach(func() {
			profileCatalog = catalog.New()
			profileCatalog.Update("foo", profilesv1.ProfileDescription{Name: "nginx-1", Description: "nginx 1"})
			catalogAPI = api.New(profileCatalog)
		})

		When("a matching profile exists", func() {
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
	})

	Context("/profiles/catalog/profile-name", func() {
		var (
			catalogName, profileName string
		)

		BeforeEach(func() {
			catalogName, profileName = "catalog", "nginx-1"
			profileCatalog = catalog.New()
			profileCatalog.Update(catalogName, profilesv1.ProfileDescription{Name: profileName, Description: "nginx 1"})
			catalogAPI = api.New(profileCatalog)
		})

		When("the requested profile exists", func() {
			It("returns the profile summary from the catalog", func() {
				req, err := http.NewRequest("GET", "/profiles", nil)
				req = mux.SetURLVars(req, map[string]string{"catalog": catalogName, "profile": profileName})
				Expect(err).NotTo(HaveOccurred())

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(catalogAPI.ProfileHandler)

				handler.ServeHTTP(rr, req)

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Body.String()).To(ContainSubstring(`{"name":"nginx-1","description":"nginx 1","catalog":"catalog"}`))
			})
		})
	})
})
