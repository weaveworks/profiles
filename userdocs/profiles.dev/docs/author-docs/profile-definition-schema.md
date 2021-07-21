---
sidebar_position: 8
---

# Profile definition schema

<!-- TODO: autogen this in like yaml or smth -->

```go
// ProfileDefinitionSpec defines the desired state of ProfileDefinition
type ProfileDefinitionSpec struct {
	ProfileDescription `json:",inline"`
	// Artifacts is a list of Profile artifacts
	Artifacts []Artifact `json:"artifacts,omitempty"`
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

// Artifact defines a bundled resource of the components for this profile.
type Artifact struct {
	// Name is the name of the Artifact
	Name string `json:"name,omitempty"`
	// Chart defines properties to access a remote chart.
	// This is an optional value. It is ignored in case Path is defined.
	// +optional
	Chart *Chart `json:"chart,omitempty"`
	// Profiles defines properties to access a remote profile.
	// +optional
	Profile *Profile `json:"profile,omitempty"`
	// Kustomize defines properties to for a kustmize artifact.
	// +optional
	Kustomize *Kustomize `json:"kustomize,omitempty"`
}

type Kustomize struct {
	// Path is the local path to the Artifact in the Profile repo.
	Path string `json:"path,omitempty"`
}

// Chart defines properties to access remote helm charts.
type Chart struct {
	// URL is the URL of the Helm repository containing a Helm chart and possible values
	// +optional
	URL string `json:"url,omitempty"`
	// Name defines the name of the chart at the remote repository
	// +optional
	Name string `json:"name,omitempty"`
	// Version defines the version of the chart at the remote repository
	// +optional
	Version string `json:"version,omitempty"`
	// Path is the local path to the Artifact in the Profile repo.
	// This is an optional value. If defined, it takes precedence over other Chart fields.
	// +optional
	Path string `json:"path,omitempty"`
	// DefaultValues holds the default values for this Helm release Artifact.
	// These can be overridden by the user, but will otherwise apply.
	// +optional
	DefaultValues string `json:"defaultValues,omitempty"`
}

// Profile defines properties for accessing a profile
type Profile struct {
	// Source defines properties of the source of the profile
	Source *Source `json:"source,omitempty"`
}

// Source defines the location of the profile
type Source struct {
	// ProfileURL is a fully qualified URL to a profile repo
	URL string `json:"url,omitempty"`

	// Branch is the git repo branch containing the profile definition (default: main)
	// +kubebuilder:default:=main
	// +optional
	Branch string `json:"branch,omitempty"`

	// Path is the location in the git repo containing the profile definition
	// +optional
	Path string `json:"path,omitempty"`

	// Tag is the git tag containing the profile definition
	// +optional
	Tag string `json:"tag,omitempty"`
}
```
