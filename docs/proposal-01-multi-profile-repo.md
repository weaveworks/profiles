<!--
How to use this template:

- Make a copy of this file in the docs/ directory
- Set the name of the file to contain the next logical number and the name of the feature
- Fill out at least the Status, Motivation and Goals/Non-Goals fields.
- Open a PR to eksctl
- Merge early and iterate

For more tips see the Contributing docs: https://github.com/weaveworks/eksctl/blob/master/CONTRIBUTING.md#proposals
-->

# Multi-profile repos with versioning

<!--
Keep the title short, simple, and descriptive. A good
title can help communicate what the proposal is and should be considered as part of
any review.
-->

## Authors

Jake Klein @aclevername

## Status

WIP


<!--
The headings here are just starting points, add more as makes sense in what you
are proposing.
-->
## Table of Contents
<!-- toc -->
- [Summary](#summary)
- [Design Details](#design-details)
- [Drawbacks (Optional)](#drawbacks-optional)
- [Open Questions / Known Unknowns](#open-questions--known-unknowns)
<!-- /toc -->

## Summary

<!--
A good summary is at least a paragraph in length and should be written with a wide audience
in mind.

This TLDR should encompass the entire document, and serve as both future documentation
and as a quick reference for people coming by to learn the proposal's purpose
without reading the entire thing.
-->
The purpose of this doc is to outline the proposed implementation of multi-profile repositories, and how to handle tagging profiles in this scenario.
Currently you are limited to 1 profile per repository and you can only reference a branch in a repository. This document proposes that 
we introduce the ability to have multiple profiles in a repository, and follow a pattern of tagging used in [kustomize](https://github.com/kubernetes-sigs/kustomize/tags)
to distinguish what tag is used for each profile.


## Design Details
### Repository Layout

An example repository `weaveworks-profiles` containing two profiles, `foo` and `bar` would be laid out as follows:
```bash
$ ls weaveworks-profiles
foo/
bar/
```

Inside `foo` and `bar` it would contain the `profiles.yaml` and any local artifacts, for example:
```bash
$ tree weaveworks-profile/foo
├── charts
│   └── nginx
│       └── Chart.yaml
├── kustomize
│   └── web-app
│       └── deployment.yaml
└── profile.yaml

```
The `profile.yaml` would contain a reference to the artifacts relative to profile directory, for example:
```yaml
$ cat foo/profile.yaml
apiVersion: profiles.fluxcd.io/v1alpha1
kind: Profile
metadata:
  name: foo
spec:
  description: Profile for deploying local nginx chart and web-app
  version: v0.0.2
  artifacts:
    - name: nginx-server
      path: charts/nginx
      kind: HelmChart
    - name: web-app
      path: kustomize/web-app
      kind: Kustomize
```

### Tagging
[kustomize](https://github.com/kubernetes-sigs/kustomize/tags) is an example of a repository that contains multiple "things" that are all versioned
independently. For example if you take a look at some of the tags it has:
```
kyaml/v0.10.19
kyaml/v0.10.18
api/v0.8.9
api/v0.8.8
cmd/config/v0.9.11
cmd/config/v0.9.10
```

The prefix before the semver version corresponds to the directory in which the product lives, with the last directory also being the name of the component. For example the repository contains:
```
tree kustomize/
├── api
    └── main.go
    ...
├── cmd
│   └── config
	└── main.go
	...
└── kyaml
    └── main.go
    ...
```

If you wanted to get version `v0.8.9` of `api` you would checkout to tag `api/v0.8.9` and inside the `api/` directory it would contain the desired code.
The case of `cmd/config` is slightly less clean, its easier to write automation around detecting new profiles if the profile directory is always at the top
level directory. Support for sub-directories is in the [Open Questions / Known Unknowns](#open-questions--known-unknowns) section below.

This approach could be used for versioning profiles in a repository.

### Profile Subscription modifications
Currently we only support referencing a repository as follows:
```
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: nginx-profile-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/nginx-profile
  branch: main
```

We want to maintain support for `branch` workflows as its useful for development, while also introducing a new `tags` approach. There are two possible approaches:

#### Approach 1

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: foo-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/weaveworks-profiles
  tag: v0.1.2
  profile: foo
```

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: bar-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/weaveworks-profiles
  branch: main
  profile: bar
```

This approach hides the fact that the git tag contains the profile name, and leaves it up to the profiles controller to concatenate the `profile` and `tag` value together.
This introduces a common `profile` (or `profileName`) field that is shared across the two for us to know which profile in the repository we are referencing and its directory (they must be equal).

#### Approach 2

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: foo-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/weaveworks-profiles
  tag: foo/v0.1.2
```

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: bar-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/weaveworks-profiles
  branch: main
  path: bar/
```

This approach makes it clear upon inspection what tag is actually being referenced in the repository, but is less clear what directory/profile inside the repository is being referenced.
It also maintains two completely different approach for using `branch` vs `tag`

## Drawbacks (Optional)

- This approach to tagging repositories that contain multiple "things" is not widely used, users will likely not be familiar with it

## Open Questions / Known Unknowns

Support for specifying sub-directories in the tag allow more granular repository layouts, but make parsing tags less clean. For example if my repository is as follows:

```bash
$ ls weaveworks-profiles
foo/
bar/
```

```
bar/v0.10.19
bar/v0.10.18
foo/v0.1.1
foo/v0.1.0
```
Its clear that it has two profiles `foo` and `bar`, and that they have there respective tags. If we added support for sub-directories we introduce more complexity, for example:

```bash
$ ls weaveworks-profiles
foo/
bar/
sub/
  path/
    foo/
```

```
bar/v0.10.19
bar/v0.10.18
foo/v0.1.1
foo/v0.1.0
sub/path/foo/v1.1.0
sub/path/foo/v1.0.0
```

Here we have two different profiles, both called `foo` which are only distinguishable by their paths. If we were to write automation to parse profile repositories we might find it awkward
to handle such scenarios. It would also result in more obscure subscription definitions, for example to adapt approach 1 would look like:
```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: foo-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/weaveworks-profiles
  tag: v0.1.2
  profile: foo
  path: sub/path
```

or

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileSubscription
metadata:
  name: foo-test
  namespace: default
spec:
  profileURL: https://github.com/weaveworks/weaveworks-profiles
  tag: v0.1.2
  profile: sub/path/foo
```
