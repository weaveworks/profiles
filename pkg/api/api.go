package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/weaveworks/profiles/pkg/catalog"
)

// API defines a catalog router.
type API struct {
	*mux.Router
	profileCatalog catalog.Catalog
}

// New returns a new mux based api router.
func New(profileCatalog catalog.Catalog) API {
	r := mux.NewRouter()
	a := API{
		Router:         r,
		profileCatalog: profileCatalog,
	}

	r.HandleFunc("/profiles", a.ProfilesHandler)
	r.HandleFunc("/profiles/{catalog}/{profile}", a.ProfileHandler)
	r.HandleFunc("/profiles/{catalog}", a.DeleteCatalogHandler).
		Methods(http.MethodDelete)

	return a
}

// ProfilesHandler is the handler for /profiles requests.
func (a *API) ProfilesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	marshalResponse(w, a.profileCatalog.Search(query))
}

// ProfileHandler is the handler for /profiles/{catalog}/{profile} requests.
func (a *API) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catalogName, profileName := vars["catalog"], vars["profile"]
	profile, found := a.profileCatalog.Get(catalogName, profileName)
	if !found {
		w.WriteHeader(404)
		return
	}
	marshalResponse(w, profile)
}

// DeleteCatalogHandler deletes a catalog
func (a *API) DeleteCatalogHandler(w http.ResponseWriter, r *http.Request) {
	a.profileCatalog.Remove(mux.Vars(r)["catalog"])
}

func marshalResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(500)
		log.Printf("failed to encode response: %s", err)
	}
}
