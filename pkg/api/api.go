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

func (a *API) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	profileName := mux.Vars(r)["name"]
	out, err := json.Marshal(a.profileCatalog.Show(profileName))
	if err != nil {
		panic(err)
	}
	_, err = w.Write(out)
	if err != nil {
		panic(err)
	}
}
