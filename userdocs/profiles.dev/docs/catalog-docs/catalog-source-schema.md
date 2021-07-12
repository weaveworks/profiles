---
sidebar_position: 2
---

# Catalog source schema

<!-- TODO: autogen this in like yaml or smth -->

```go
// ProfileCatalogSourceSpec defines the desired state of ProfileCatalogSource
type ProfileCatalogSourceSpec struct {
	// Profiles is the list of profiles exposed by the catalog
	// +optional
	Profiles []ProfileCatalogEntry `json:"profiles,omitempty"`
	// Repos contains a list of repositories to scan for profiles
	// +optional
	Repos []Repository `json:"repositories,omitempty"`
}

// Repository defines the list of repositories to scan for profiles
type Repository struct {
	// URL is the URL of the repository. When using SSH credentials to access
	// must be in format ssh://git@github.com/weaveworks/profiles-examples
	// When using username/password must be in format
	// https://github.com/weaveworks/profiles-examples
	URL string `json:"url,omitempty"`
	// The secret name containing the Git credentials.
	// For HTTPS repositories the secret must contain username and password
	// fields.
	// For SSH repositories the secret must contain identity, identity.pub and
	// known_hosts fields.
	// +optional
	SecretRef *meta.LocalObjectReference `json:"secretRef,omitempty"`
}

// ProfileCatalogEntry defines details about a given profile.
type ProfileCatalogEntry struct {
	// Tag
	// +optional
	// +kubebuilder:validation:Pattern=^([a-zA-Z\-]+\/)?(v)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$

	Tag string `json:"tag,omitempty"`
	// CatalogSource is the name of the catalog the profile is listed in
	// +optional
	CatalogSource string `json:"catalog,omitempty"`
	// URL is the full URL path to the profile.yaml
	// +optional
	URL                string `json:"url,omitempty"`
	Name               string `json:"name,omitempty"`
	ProfileDescription `json:",inline"`
}

// ProfileDescription defines details about a given profile.
type ProfileDescription struct {
	// Profile description
	Description string `json:"description,omitempty"`
	// Maintainer is the name of the author(s)
	// +optional
	Maintainer string `json:"maintainer,omitempty"`
	// Prerequisites are a list of dependencies required by the profile
	// +optional
	Prerequisites []string `json:"prerequisites,omitempty"`
}
```
