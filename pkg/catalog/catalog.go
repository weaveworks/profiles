package catalog

import (
	"strings"

	"github.com/weaveworks/profiles/api/v1alpha1"
)

type Catalog struct {
	profiles []v1alpha1.ProfileDescription
}

func New() *Catalog {
	return &Catalog{profiles: []v1alpha1.ProfileDescription{}}
}

func (c *Catalog) Add(p ...v1alpha1.ProfileDescription) {
	c.profiles = append(c.profiles, p...)
}

func (c *Catalog) Search(name string) []v1alpha1.ProfileDescription {
	var profiles []v1alpha1.ProfileDescription
	for _, p := range c.profiles {
		if strings.Contains(p.Name, name) {
			profiles = append(profiles, p)
		}
	}
	return profiles
}
