package api_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/api"
	"github.com/weaveworks/profiles/pkg/catalog"
)

var _ = Describe("Api", func() {
	var (
		catalogAPI     api.API
		profileCatalog *catalog.Catalog
	)
	Context("/profiles", func() {
		BeforeEach(func() {
			profileCatalog = catalog.New(logr.Discard())
			profileCatalog.Add(v1alpha1.ProfileDescription{Name: "nginx-1", Description: "nginx 1"})
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
				Expect(rr.Body.String()).To(Equal(`[{"name":"nginx-1","description":"nginx 1"}]`))
			})
		})
	})

})
