package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/weaveworks/profiles/pkg/catalog"
)

// API defines a catalog router.
type API struct {
	*mux.Router
	profileCatalog *catalog.Catalog
}

// New creates a new mux based api router.
func New(profileCatalog *catalog.Catalog) API {
	r := mux.NewRouter()
	a := API{
		Router:         r,
		profileCatalog: profileCatalog,
	}

	r.HandleFunc("/profiles", a.ProfilesHandler)
	return a
}

// ProfilesHandler is the handler for /profiles events.
func (a *API) ProfilesHandler(w http.ResponseWriter, r *http.Request) {
	profileName := r.URL.Query().Get("name")
	out, err := json.Marshal(a.profileCatalog.Search(profileName))
	if err != nil {
		panic(err)
	}
	_, err = w.Write(out)
	if err != nil {
		panic(err)
	}
}
