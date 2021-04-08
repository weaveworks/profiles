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
	profileCatalog *catalog.Catalog
}

// New returns a new mux based api router.
func New(profileCatalog *catalog.Catalog) API {
	r := mux.NewRouter()
	a := API{
		Router:         r,
		profileCatalog: profileCatalog,
	}

	r.HandleFunc("/profiles", a.ProfilesHandler)
	r.HandleFunc("/profiles/{name}", a.ProfileHandler)

	return a
}

// ProfilesHandler is the handler for /profiles requests.
func (a *API) ProfilesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	marshalResponse(w, a.profileCatalog.Search(query))
}

func (a *API) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	profileName := mux.Vars(r)["name"]
	marshalResponse(w, a.profileCatalog.Get(profileName))
}

func marshalResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode response: %s", err)
	}
}
