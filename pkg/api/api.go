package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

//go:generate counterfeiter -o fakes/fake_catalog.go . Catalog
// Catalog is an interface for the Catalog
type Catalog interface {
	// Get will return a specific profile from the catalog
	Get(sourceName, profileName string) *profilesv1.ProfileDescription
	// GetWithVersion will return a specific profile from the catalog
	GetWithVersion(sourceName, profileName, version string) *profilesv1.ProfileDescription
	// GetGreaterThan returns all profiles which are of a greater version for a given profile with a version.
	ProfilesGreaterThanVersion(sourceName, profileName, version string) []profilesv1.ProfileDescription
	// Search will return a list of profiles which match query
	Search(query string) []profilesv1.ProfileDescription
}

// API defines a catalog router.
type API struct {
	*mux.Router
	profileCatalog Catalog
}

// New returns a new mux based api router.
func New(profileCatalog Catalog) API {
	r := mux.NewRouter()
	a := API{
		Router:         r,
		profileCatalog: profileCatalog,
	}

	r.HandleFunc("/profiles", a.ProfilesHandler)
	r.HandleFunc("/profiles/{catalog}/{profile}", a.ProfileHandler)
	r.HandleFunc("/profiles/{catalog}/{profile}/{version}", a.ProfileWithVersionHandler)
	r.HandleFunc("/profiles/{catalog}/{profile}/{version}/available_updates", a.ProfileGreaterThanVersionHandler)

	return a
}

// ProfilesHandler is the handler for /profiles requests.
func (a *API) ProfilesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	marshalResponse(w, a.profileCatalog.Search(query))
}

// ProfileHandler is the handler for /profiles/{catalog}/{profile} requests.
func (a *API) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	sourceName, profileName := mux.Vars(r)["catalog"], mux.Vars(r)["profile"]
	result := a.profileCatalog.Get(sourceName, profileName)
	if result == nil {
		w.WriteHeader(404)
		return
	}
	marshalResponse(w, result)
}

// ProfileWithVersionHandler is the handler for /profiles/{catalog}/{profile}/{version} requests.
func (a *API) ProfileWithVersionHandler(w http.ResponseWriter, r *http.Request) {
	sourceName, profileName, catalogVersion := mux.Vars(r)["catalog"], mux.Vars(r)["profile"], mux.Vars(r)["version"]
	result := a.profileCatalog.GetWithVersion(sourceName, profileName, catalogVersion)
	if result == nil {
		w.WriteHeader(404)
		return
	}
	marshalResponse(w, result)
}

// ProfileGreaterThanVersionHandler is the handler for /profiles/{catalog}/{profile}/{version} requests.
func (a *API) ProfileGreaterThanVersionHandler(w http.ResponseWriter, r *http.Request) {
	sourceName, profileName, catalogVersion := mux.Vars(r)["catalog"], mux.Vars(r)["profile"], mux.Vars(r)["version"]
	result := a.profileCatalog.ProfilesGreaterThanVersion(sourceName, profileName, catalogVersion)
	if result == nil {
		w.WriteHeader(404)
		return
	}
	marshalResponse(w, result)
}

func marshalResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(500)
		log.Printf("failed to encode response: %s", err)
	}
}
