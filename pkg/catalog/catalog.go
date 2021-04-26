package catalog

import (
	"strings"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// Catalog provides an in-memory cache of profiles from the cluster which can be queried easily.
type Catalog struct {
	profiles []profilesv1.ProfileDescription
}

// New creates a new, empty catalog.
func New() *Catalog {
	return &Catalog{profiles: []profilesv1.ProfileDescription{}}
}

// Add adds p profiles to the catalog.
func (c *Catalog) Add(catalogName string, profiles ...profilesv1.ProfileDescription) {
	for _, p := range profiles {
		p.CatalogSource = catalogName
		c.profiles = append(c.profiles, p)
	}
}

// Search returns profile descriptions that contain `name` in their names.
func (c *Catalog) Search(name string) []profilesv1.ProfileDescription {
	var profiles []profilesv1.ProfileDescription
	for _, p := range c.profiles {
		if strings.Contains(p.Name, name) {
			profiles = append(profiles, p)
		}
	}

	return profiles
}

// Get returns the profile description `profileName`.
func (c *Catalog) Get(catalogName, profileName string) profilesv1.ProfileDescription {
	for _, p := range c.profiles {
		if p.Name == profileName && p.CatalogSource == catalogName {
			return p
		}
	}

	return profilesv1.ProfileDescription{}
}
