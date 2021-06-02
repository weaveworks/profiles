# profiles
Gitops native package management.

<!--
To update the TOC, install https://github.com/kubernetes-sigs/mdtoc
and run: mdtoc -inplace README.md
-->

 <!-- toc -->
- [Getting started](#getting-started)
  - [Local environment using <a href="https://kind.sigs.k8s.io/">Kind</a>](#local-environment-using-kind)
  - [Installing Profiles](#installing-profiles)
- [Development](#development)
  - [Tests](#tests)
  - [Release process](#release-process)
  - [Dev Tags](#dev-tags)
- [Terminology](#terminology)
  - [Profile](#profile)
  - [Catalog](#catalog)
  - [Profile Catalog Source](#profile-catalog-source)
  - [Profile Catalog Source Controller](#profile-catalog-source-controller)
- [Current Architecture](#current-architecture)
  - [Catalogs and Sources](#catalogs-and-sources)
  - [Profile Installation](#profile-installation)
- [Roadmap](#roadmap)
  - [Profiles](#profiles)
  - [Catalogs](#catalogs)
<!-- /toc -->

## Getting started

### Local environment using [Kind](https://kind.sigs.k8s.io/)

1. Set up local environment: `make local-env`.

    This will start a local `kind` cluster and installs
    the `profiles` and `flux` components.

1. Deploy an example catalog source `kubectl apply -f examples/profile-catalog-source.yaml`

### Installing Profiles

1. Profiles can be installed using [pctl](https://github.com/weaveworks/pctl).

## Development

### Tests

1. All tests can be run with `make test`.

1. Acceptance tests can be run with `make acceptance`.

1. For further commands, run `make help`.

### Release process
There are some manual steps right now, should be streamlined soon.

Steps:

1. Create a new release notes file:
	```sh
	touch docs/release_notes/<version>.md
	```

1. Copy-and paste the release notes from the draft on the releases page into this file.
    _Note: sometimes the release drafter is a bit of a pain, verify that the notes are
    correct by doing something like: `git log --first-parent tag1..tag2`._

1. PR the release notes into main.

1. Create and push a tag with the new version:
	```sh
	git tag <version>
	git push origin <version>
	```

1. The `Create release` action should run. Verify that:
	1. The release has been created in Github
		1. With the correct assets
		1. With the correct release notes
	1. The image has been pushed to docker
	1. The image can be pulled and used in a deployment
    
### Dev Tags

As part of pushing a new branch to profiles, a new dev tag will be created for that branch. This action is part of
the procedure around working with profiles and [pctl](https://github.com/weaveworks/pctl), it's companion CLI tool. The details of this are
explained in `pctl`'s README section [Working with Profiles](https://github.com/weaveworks/pctl#working-with-profiles).

The dev tag's format is as follows: `<currentLatestReleaseTag>-<branch-name>`.

## Terminology

### Profile

A Profile is a "package" of Kubernetes deployable objects, known as Artifacts, and configurable values.
Artifacts are one of: Helm Chart; Helm Release; raw yaml; Kustomize patch; Profile (nested).

For an example, see the [profiles-examples](https://github.com/weaveworks/profiles-examples).

### Catalog

A Catalog is an in-memory cache of Profiles. There is one Catalog per running [Profile Controller](#profile-controller).
The Catalog is queryable via [pctl](https://github.com/weaveworks/pctl) or the API directly which runs alongside the Profiles Controller.

Profiles can be added to the Catalog by creating a [`ProfileCatalogSource`](#profile-catalog-source).

### Profile Catalog Source

A `ProfileCatalogSource` is a custom resource through which approved Profiles can be managed in the [Catalog](#catalog)

```go
// ProfileCatalogSourceSpec defines the desired state of ProfileCatalogSource
type ProfileCatalogSourceSpec struct {
	// Profiles is the list of profiles exposed by the catalog
	Profiles []ProfileDescription `json:"profiles,omitempty"`
}

// ProfileDescription defines details about a given profile.
type ProfileDescription struct {
	// Profile name
	Name string `json:"name,omitempty"`
	// Profile description
	Description string `json:"description,omitempty"`
	// Version
	// +optional
	Version string `json:"version,omitempty"`
	// CatalogSource is the name of the catalog the profile is listed in
	// +optional
	CatalogSource string `json:"catalog,omitempty"`
	// URL is the full URL path to the profile.yaml
	// +optional
	URL string `json:"url,omitempty"`
	// Maintainer is the name of the author(s)
	// +optional
	Maintainer string `json:"maintainer,omitempty"`
	// Prerequisites are a list of dependencies required by the profile
	// +optional
	Prerequisites []string `json:"prerequisites,omitempty"`
}
```

Profiles can therefore be grouped and namespaced within the Catalog.


### Profile Catalog Source Controller

The Profile Catalog Source Controller reconciles `ProfileCatalogSource` resources.
See architecture diagrams below for what the reconciliation process does.

## Current Architecture

### Catalogs and Sources

<!--
To update this diagram go to https://miro.com/app/board/o9J_lI2seIg=/
edit, export, save as image (size small) and commit. Easy.
-->
Illustration of how Profiles are added to the Catalog and how they can then be queried via the Catalog API:

![](/docs/assets/catalog.jpg)

### Profile Installation

To see how a profile installation works, take a look at the [pctl documentation](https://github.com/weaveworks/pctl)

## Roadmap

### Profiles

Install:
- [x] Install a simple profile which contains a single Helm release artifact
- [x] Install a simple profile which contains a raw yaml artifact (k8s object manifest)
- [x] Install a simple profile which contains another profile (single nesting)
- [x] Install a profile which contains a mix of all artifact types
- [x] Install a profile which contains nested profiles to depth N
- [x] Install a profile with `pctl` in a gitops way (ie there is a PR involved, and nobody touches the cluster)
- [x] Install a profile which is listed in the catalog
- [x] Install a profile which is NOT listed in the catalog
- [ ] Install a private profile

Configure:
- [x] Configure a Helm release artifact installation
- [ ] Apply Kustomise patches
- [ ] Configure different values across multiple artifacts

Update:
- [ ] Discover when there is a newer version available
- [ ] Update a profile

### Catalogs

Catalog sources:
- [x] Create a catalog source
- [x] Delete a catalog source
- [ ] Grant/Revoke access to CatalogSources

Catalog management:
- [x] Add profiles to the catalog
- [X] Update profiles in the catalog
- [ ] Delete profiles from the catalog

API:
- [x] Search for profiles in the catalog
- [x] Get more information about a profile in the catalog
