package catalog

import (
	"fmt"
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
func (c *Catalog) Add(p ...profilesv1.ProfileDescription) {
	c.profiles = append(c.profiles, p...)

	added := map[string]struct{}{}
	for i, p := range c.profiles {
		if _, ok := added[p.Name]; ok {
			c.log.Info(fmt.Sprintf("profile %s already exists in catalog", p.Name))

			// We add all profiles at the top of the func then remove dupes here because
			// if multiples are being added then it is less straightforward to check against
			// what is already in the catalog as well as dupes in the Add list
			ret := c.profiles[:i]
			c.profiles = append(ret, c.profiles[i+1:]...)

			continue
		}
		added[p.Name] = struct{}{}
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
