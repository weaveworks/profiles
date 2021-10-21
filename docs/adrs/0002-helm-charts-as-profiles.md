# 2. Helm Charts as Profiles

Date: 2021-10-20

## Status

Accepted

## Context

We want to build a set of packages for Kubernetes, that we can use as building
blocks for configuring clusters.

These are versioned, and can provide simplified configuration for these
packages, where they might incorporate sub-packages.

## Decision

Profiles are Helm charts.

They can be Helm charts with [Subcharts](https://helm.sh/docs/chart_template_guide/subcharts_and_globals/) or [dependencies](https://helm.sh/docs/chart_best_practices/dependencies/#helm).

There's a conceptual difference between Profiles and Helm charts, but the
technical implementation is implemented entirely in Helm charts.

To identify whether or not a Helm chart is a Profile:

This proposes `weave.works/profile: name-of-profile` as the annotation, for
example `Chart.yaml`.

```yaml
name: demo-profile
version: 0.0.1
annotations:
  weave.works/profile: "A Demo Profile"
```

In the UI, these will be referred to as Profiles, and the tooling for Profiles
will be updated to reflect these changes.

"Installation" of a profile, means writing a `HelmRelease` object to the gitops
configuration repository. It's assumed that the appropriate `HelmRepository` resource already exists in cluster.

## Consequences

 * Harder to do dependsOn from Flux HelmReleases
 * Easier to debug[^debug], they are just Helm charts
 * Easier transport mechanisms[^transport], again, they're just Helm charts
 * Lighter to install, installing a profile doesn't require cloning upstream
   repositories
 * Consumes less resources in the cluster
 * Piggybacks on the extensive documentation for Helm

[^debug]: The current approach has four layers, Profiles, Flux Helm (Releases
  and Charts), Helm and the Helm templated resources, this removes one of the
  layers which should simplify debugging.

[^transport]: There is a well-defined Helm API. and `helm` can access Charts
  using several different protocols (http/https/sftp).
