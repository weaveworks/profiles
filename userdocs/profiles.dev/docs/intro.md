---
sidebar_position: 1
---

# Introduction

Welcome to Profiles: _The_ GitOps Native Package Manager :sparkles:

:::danger _(highway to the)_ Danger Zone
Profiles is still highly experimental and likely to experience sudden and significant
API change as we continue to develop based on incoming feedback.

Please always refer to the docs and release notes, and alert us to any issues
or inconsistencies.

Thank you!
:::

## What are Profiles?

In short, a profile is a Kubernetes package. Much as you would expect from another package
manager (`brew`, `apt`, etc), the Profiles mechanism allows Kubernetes operators
to manage their clusters in the same way they manage their host systems. Available
packages are listed somewhere searchable; they are versioned; they are tested and verified;
they are installable by way of a simple command; they are updatable and upgradable.
Moreover they are reliable and uniformly manageable through a simple intuitive cli tool.

Sticking with the `apt` analogy, let's look at one of the most basic and ubiquitous
packages: `coreutils`. By installing `coreutils`, users understand that they are
installing a bulk set of the most commonly used and useful tools on their unix system.
They don't care much what the full roster is, just that their system is now more operable
than it was before; that the various pieces work; that the whole thing can be updated
or removed just as simply.

Profiles provides the same service for cluster operators. On spinning up a new cluster or
fleet of clusters, operators can then apply an approved set of profiles to get things
to a uniform, operable standard. They could choose to install, say, a standard Observability
package, maybe a Logging set too. If they have some in-house collections, they are able to create
their own custom profiles to bundle these up, and apply them across their infrastructure.
They even have the option of creating a single "install all the things" profile to bring
their most commonly used collections under one reliable and reusable package, giving them just one
thing to do to bootstrap a new cluster to their org specifications.

Profiles is designed from its core to follow the best practices
of GitOps, that is; managing infrastructure as code. Profiles therefore uses [Flux](https://fluxcd.io/)
to fit into users' existing cluster management practices.

:::tip Heard the word "profiles" too much?
No worries, here's the breakdown:

- A profile: an individual package of Kubernetes components. Lives in a git repo "upstream"
  of users who have installed it on their cluster
- `N` profiles: more than one profile
- Profiles (capital 'P'): a blanket term for this concept of a GitOps native package management mechanism
- Pctl: the CLI tool. Use this to install and manage profiles on your cluster
:::

---------------------

## Core concepts & Terminology

### GitOps

GitOps is a way of managing your infrastructure and applications so that the whole
system is described declaratively and version controlled (most likely in a Git repository),
and having an automated process that ensures that the deployed environment matches the state specified in a repository.

For more information, take a look at [“What is GitOps?"](https://www.gitops.tech/#what-is-gitops).

The GitOps tool leveraged by Profiles is Flux. Please refer to [their documentation](https://fluxcd.io/) for more
information.

### Profile

A Profile is a "package" of Kubernetes deployable objects, known as Artifacts, along with any configurable values.
An artifact can be one of:
- Helm Chart
- Raw Kubernetes yaml
- Kustomize patch
- Another nested Profile

See the [example repo](https://github.com/weaveworks/profiles-examples) for a full collection of working
profile artifact types.

### pctl

Profiles are installed and managed via the official CLI: `pctl`.

### Profile repository

A profile repository holds the definition and other necessary items for an upstream profile
or profiles. A single repository can hold multiple profiles, separated by directories.

### GitOps Repository

A GitOps repository is one which is synced to your cluster via Flux. The manifests of all
components applied to your cluster should be managed through this repo.

When installing profiles, the generated manifests will be saved to this repo.

Refer to the [Flux documentation](https://fluxcd.io/) for more detail on how to work with Flux.

### Install

To install a Profile means to generate the required component manifests for a Profile and commit
them to a GitOps repo. From there, Flux will recognise the changes and apply the manifests.

For more information, see the [documentation](/docs/installer-docs/installing-via-gitops).

### Catalog

A Catalog is an in-cluster cache of Profile references. The mechanism allows cluster admins to define
a list of profiles which may be applied to a cluster, and also provides users with another
method of installation.

Users can also see when newer versions of profiles are published, as well as query the cache
for more profiles which may meet their needs.

There is one Catalog per running [Profile Catalog Source Controller](#profile-controller).
The Catalog is queryable via [pctl](https://github.com/weaveworks/pctl).

Profiles can be added to the Catalog by creating a [`ProfileCatalogSource`](#profile-catalog-source).

### Profile Catalog Source

A `ProfileCatalogSource` is a custom resource through which approved profiles can be managed in the [Catalog](#catalog).

### Profile Catalog Source Controller

The Profile Catalog Source Controller reconciles `ProfileCatalogSource` resources.
See architecture diagrams below for what the reconciliation process does.

 
