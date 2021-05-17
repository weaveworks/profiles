package catalog

import (
	"strings"
	"sync"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// Catalog provides an in-memory cache of profiles from the cluster which can be queried easily.
//type Catalog map[string][]profilesv1.ProfileDescription
type Catalog struct {
	m sync.Map
}

// New creates a new, empty catalog.
func New() *Catalog {
	return &Catalog{
		m: sync.Map{},
	}
}

// Update updates the catalog by replacing existing profiles with new profiles
func (c *Catalog) Update(sourceName string, profiles ...profilesv1.ProfileDescription) {
	for i := range profiles {
		profiles[i].CatalogSource = sourceName
	}
	c.m.Store(sourceName, profiles)
}

// Remove removes the specified catalog.
func (c *Catalog) Remove(sourceName string) {
	c.m.Delete(sourceName)
}

// Search returns profile descriptions that contain `name` in their names.
func (c *Catalog) Search(name string) []profilesv1.ProfileDescription {
	var ret []profilesv1.ProfileDescription
	c.m.Range(func(key, value interface{}) bool {
		for _, p := range value.([]profilesv1.ProfileDescription) {
			if strings.Contains(p.Name, name) {
				ret = append(ret, p)
			}
		}
		return true
	})
	return ret
}

// Get returns the profile description `profileName`.
func (c *Catalog) Get(sourceName, profileName string) *profilesv1.ProfileDescription {
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}
	for _, p := range profiles.([]profilesv1.ProfileDescription) {
		if p.Name == profileName && p.CatalogSource == sourceName {
			return p.DeepCopy()
		}
	}
	return nil
}

// GetWithVersion returns the profile description `profileName` with the given version.
func (c *Catalog) GetWithVersion(sourceName, profileName, version string) *profilesv1.ProfileDescription {
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}
	for _, p := range profiles.([]profilesv1.ProfileDescription) {
		if p.Name == profileName && p.CatalogSource == sourceName && p.Version == version {
			return p.DeepCopy()
		}
	}
	return nil
}
