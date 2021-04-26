package catalog

import (
	"strings"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// Catalog provides an in-memory cache of profiles from the cluster which can be queried easily.
type Catalog map[string][]profilesv1.ProfileDescription

// New creates a new, empty catalog.
func New() Catalog {
	return map[string][]profilesv1.ProfileDescription{}
}

// Update updates the catalog by replacing existing profiles with new profiles
func (c Catalog) Update(sourceName string, profiles ...profilesv1.ProfileDescription) {
	for i := range profiles {
		profiles[i].CatalogSource = sourceName
	}
	c[sourceName] = profiles
}

// Remove removes the specified catalog.
func (c Catalog) Remove(sourceName string) {
	delete(c, sourceName)
}

// Search returns profile descriptions that contain `name` in their names.
func (c Catalog) Search(name string) []profilesv1.ProfileDescription {
	var ret []profilesv1.ProfileDescription
	for _, profiles := range c {
		for _, p := range profiles {
			if strings.Contains(p.Name, name) {
				ret = append(ret, p)
			}
		}
	}

	return ret
}

// Get returns the profile description `profileName`.
func (c Catalog) Get(sourceName, profileName string) *profilesv1.ProfileDescription {
	profiles := c[sourceName]
	for _, p := range profiles {
		if p.Name == profileName && p.CatalogSource == sourceName {
			return p.DeepCopy()
		}
	}

	return nil
}
