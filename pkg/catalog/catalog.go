package catalog

import (
	"sort"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/fluxcd/pkg/version"

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
func (c *Catalog) GetWithVersion(sourceName, profileName, profileVersion string) *profilesv1.ProfileDescription {
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}

	if profileVersion == "latest" {
		versions := c.GetGreaterThan(sourceName, profileName, profileVersion)
		if len(versions) == 0 {
			return nil
		}
		return versions[0].DeepCopy()
	}

	for _, p := range profiles.([]profilesv1.ProfileDescription) {
		if p.Name == profileName && p.CatalogSource == sourceName && p.Version == profileVersion {
			return p.DeepCopy()
		}
	}
	return nil
}

type profileDescriptionWithVersion struct {
	profileDescription profilesv1.ProfileDescription
	semverVersion      *semver.Version
}

// GetGreaterThan returns all profiles which are of a greater version for a given profile with a version.
func (c *Catalog) GetGreaterThan(sourceName, profileName, profileVersion string) []profilesv1.ProfileDescription {
	var profilesWithValidVersion []profileDescriptionWithVersion
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}
	cv, err := version.ParseVersion(profileVersion)
	if err != nil && profileVersion != "latest" {
		return nil
	}
	for _, p := range profiles.([]profilesv1.ProfileDescription) {
		v, err := version.ParseVersion(p.Version)
		if err != nil {
			continue
		}
		if p.Name == profileName {
			if profileVersion == "latest" || v.GreaterThan(cv) {
				profilesWithValidVersion = append(profilesWithValidVersion, profileDescriptionWithVersion{profileDescription: p, semverVersion: v})
			}
		}
	}

	if len(profilesWithValidVersion) == 0 {
		return nil
	}

	sort.SliceStable(profilesWithValidVersion, func(i, j int) bool {
		return profilesWithValidVersion[j].semverVersion.LessThan(profilesWithValidVersion[i].semverVersion)
	})
	var result []profilesv1.ProfileDescription
	for _, p := range profilesWithValidVersion {
		result = append(result, p.profileDescription)
	}
	return result
}
