package catalog

import (
	"strings"

	"github.com/weaveworks/profiles/api/v1alpha1"
)

// Catalog contains profiles which were gathered from the underlying cluster.
type Catalog struct {
	profiles []v1alpha1.ProfileDescription
}

// New creates a new, empty catalog.
func New() *Catalog {
	return &Catalog{profiles: []v1alpha1.ProfileDescription{}}
}

// Add adds p profiles to the catalog.
func (c *Catalog) Add(p ...v1alpha1.ProfileDescription) {
	c.profiles = append(c.profiles, p...)
}

// Search returns catalogs which contain the given `name`.
func (c *Catalog) Search(name string) []v1alpha1.ProfileDescription {
	var profiles []v1alpha1.ProfileDescription
	for _, p := range c.profiles {
		if strings.Contains(p.Name, name) {
			profiles = append(profiles, p)
		}
	}
	return profiles
}
