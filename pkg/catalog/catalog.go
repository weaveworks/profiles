package catalog

import (
	"sort"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/fluxcd/pkg/version"
	"github.com/go-logr/logr"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// Catalog provides an in-memory cache of profiles from the cluster which can be queried easily.
//type Catalog map[string][]profilesv1.ProfileCatalogEntry
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
func (c *Catalog) Update(sourceName string, profiles ...profilesv1.ProfileCatalogEntry) {
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
func (c *Catalog) Search(name string) []profilesv1.ProfileCatalogEntry {
	var ret []profilesv1.ProfileCatalogEntry
	c.m.Range(func(key, value interface{}) bool {
		for _, p := range value.([]profilesv1.ProfileCatalogEntry) {
			if strings.Contains(p.Name, name) {
				ret = append(ret, p)
			}
		}
		return true
	})
	return ret
}

// Search returns `all` profile descriptions.
func (c *Catalog) SearchAll() []profilesv1.ProfileCatalogEntry {
	var ret []profilesv1.ProfileCatalogEntry
	c.m.Range(func(key, value interface{}) bool {
		for _, p := range value.([]profilesv1.ProfileCatalogEntry) {
			ret = append(ret, p)
		}
		return true
	})
	return ret
}

// Get returns the profile description `profileName`.
func (c *Catalog) Get(sourceName, profileName string) *profilesv1.ProfileCatalogEntry {
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}
	for _, p := range profiles.([]profilesv1.ProfileCatalogEntry) {
		if p.Name == profileName && p.CatalogSource == sourceName {
			return &p
		}
	}
	return nil
}

// GetWithVersion returns the profile description `profileName` with the given version.
func (c *Catalog) GetWithVersion(logger logr.Logger, sourceName, profileName, profileVersion string) *profilesv1.ProfileCatalogEntry {
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}

	if profileVersion == "latest" {
		versions := c.ProfilesGreaterThanVersion(logger, sourceName, profileName, profileVersion)
		if len(versions) == 0 {
			return nil
		}
		return &versions[0]
	}

	for _, p := range profiles.([]profilesv1.ProfileCatalogEntry) {
		if p.Name == profileName && p.CatalogSource == sourceName && profilesv1.GetVersionFromTag(p.Tag) == profileVersion {
			return &p
		}
	}
	return nil
}

type profileDescriptionWithVersion struct {
	profileDescription profilesv1.ProfileCatalogEntry
	semverVersion      *semver.Version
}

// ProfilesGreaterThanVersion returns all profiles which are of a greater version for a given profile with a version.
// If set to "latest" all versions are returned. Versions are ordered in descending order
func (c *Catalog) ProfilesGreaterThanVersion(logger logr.Logger, sourceName, profileName, profileVersion string) []profilesv1.ProfileCatalogEntry {
	var profilesWithValidVersion []profileDescriptionWithVersion
	profiles, ok := c.m.Load(sourceName)
	if !ok {
		return nil
	}
	cv, err := version.ParseVersion(profileVersion)
	if err != nil && profileVersion != "latest" {
		return nil
	}
	for _, p := range profiles.([]profilesv1.ProfileCatalogEntry) {
		v, err := version.ParseVersion(profilesv1.GetVersionFromTag(p.Tag))
		if err != nil {
			logger.Error(err, "failed to parse profile version", "profile", p)
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
	var result []profilesv1.ProfileCatalogEntry
	for _, p := range profilesWithValidVersion {
		result = append(result, p.profileDescription)
	}
	return result
}
