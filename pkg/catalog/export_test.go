package catalog

import "github.com/weaveworks/profiles/api/v1alpha1"

func (c *Catalog) Profiles() []v1alpha1.ProfileDescription {
	return c.profiles
}
