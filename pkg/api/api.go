package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

//go:generate counterfeiter -o fakes/fake_catalog.go . Catalog
// Catalog is an interface for the Catalog
type Catalog interface {
	// Get will return a specific profile from the catalog
	Get(sourceName, profileName string) *profilesv1.ProfileCatalogEntry
	// GetWithVersion will return a specific profile from the catalog
	GetWithVersion(logger logr.Logger, sourceName, profileName, version string) *profilesv1.ProfileCatalogEntry
	// ProfilesGreaterThanVersion returns all profiles which are of a greater version for a given profile with a version.
	ProfilesGreaterThanVersion(logger logr.Logger, sourceName, profileName, version string) []profilesv1.ProfileCatalogEntry
	// Search will return a list of profiles which match query
	Search(query string) []profilesv1.ProfileDescription
	// Search will return a list of all profiles 
	Search() []profilesv1.ProfileDescription
}

// API defines a catalog router.
type API struct {
	*mux.Router
	profileCatalog Catalog
	logger         logr.Logger
}

// New returns a new mux based api router.
func New(profileCatalog Catalog, logger logr.Logger) API {
	r := mux.NewRouter()
	a := API{
		Router:         r,
		profileCatalog: profileCatalog,
		logger:         logger,
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
	logger := a.logger.WithValues("endpoint", r.URL.Path, "name", query)
	if query == "" {
		result := a.profileCatalog.Search()
	}
	result := a.profileCatalog.Search(query)
	logger.Info("found profiles", "profiles", result)
	marshalResponse(w, logger, result)
}

// ProfileHandler is the handler for /profiles/{catalog}/{profile} requests.
func (a *API) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	sourceName, profileName := mux.Vars(r)["catalog"], mux.Vars(r)["profile"]
	logger := a.logger.WithValues("endpoint", r.URL.Path, "catalog", sourceName, "profile", profileName)
	if sourceName == "" || profileName == "" {
		a.logger.Error(fmt.Errorf("missing query param"), "profile and/or catalog not set")
		a.logAndWriteHeader(w, http.StatusBadRequest)
		return
	}
	result := a.profileCatalog.Get(sourceName, profileName)
	if result == nil {
		logger.Info("profile not found")
		a.logAndWriteHeader(w, http.StatusNotFound)
		return
	}
	logger.Info("profile found", "profile", result)
	marshalResponse(w, logger, result)
}

// ProfileWithVersionHandler is the handler for /profiles/{catalog}/{profile}/{version} requests.
func (a *API) ProfileWithVersionHandler(w http.ResponseWriter, r *http.Request) {
	sourceName, profileName, catalogVersion := mux.Vars(r)["catalog"], mux.Vars(r)["profile"], mux.Vars(r)["version"]
	logger := a.logger.WithValues("endpoint", r.URL.Path, "catalog", sourceName, "profile", profileName, "version", catalogVersion)
	if sourceName == "" || profileName == "" || catalogVersion == "" {
		a.logger.Error(fmt.Errorf("missing query param"), "catalog, profile and/or version not set")
		a.logAndWriteHeader(w, http.StatusBadRequest)
		return
	}
	result := a.profileCatalog.GetWithVersion(logger, sourceName, profileName, catalogVersion)
	if result == nil {
		logger.Info("profile not found")
		a.logAndWriteHeader(w, http.StatusNotFound)
		return
	}
	logger.Info("profile found", "profile", result)
	marshalResponse(w, logger, result)
}

// ProfileGreaterThanVersionHandler is the handler for /profiles/{catalog}/{profile}/{version}/available_updates requests.
func (a *API) ProfileGreaterThanVersionHandler(w http.ResponseWriter, r *http.Request) {
	sourceName, profileName, catalogVersion := mux.Vars(r)["catalog"], mux.Vars(r)["profile"], mux.Vars(r)["version"]
	logger := a.logger.WithValues("endpoint", r.URL.Path, "catalog", sourceName, "profile", profileName, "version", catalogVersion)
	if sourceName == "" || profileName == "" || catalogVersion == "" {
		a.logger.Error(fmt.Errorf("missing query param"), "catalog, profile and/or version not set")
		a.logAndWriteHeader(w, http.StatusBadRequest)
		return
	}
	result := a.profileCatalog.ProfilesGreaterThanVersion(logger, sourceName, profileName, catalogVersion)
	if len(result) == 0 {
		logger.Info("profiles not found")
		a.logAndWriteHeader(w, http.StatusNotFound)
		return
	}
	logger.Info("profile found", "profile", result)
	marshalResponse(w, logger, result)
}

func marshalResponse(w http.ResponseWriter, logger logr.Logger, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(500)
		logger.Error(err, "failed to encode response")
	}
}

func (a *API) logAndWriteHeader(w http.ResponseWriter, statusCode int) {
	a.logger.Info(fmt.Sprintf("returning %d", statusCode), "statuscode", statusCode)
	w.WriteHeader(statusCode)
}
