---
sidebar_position: 1
---

# How a profile is structured

## `profile.yaml` contents

A profile is defined in a single file which **must** be named `profile.yaml`.
This file lives at the root of the profile directory.

The following fields are required:

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileDefinition
metadata:
  name: nginx # the name of your profile 
spec:
  description: Profile for deploying local nginx chart # a brief description of what your profile installs
```

These fields will show up later when a user searches for this profile in the catalog, as shown [here](https://profiles.dev/docs/installer-docs/using-catalogs#inspecting-profiles-in-the-catalog).

These fields are optional:

```yaml
# ...
spec:
  # ...
  maintainer: weaveworks # the name(s) of the profile author
  prerequisites:
  - kubernetes 1.19 # a list of strings detailing things the profile needs to run.
  - # this field is not processed at the moment, but will be soon.
```

Finally, the `spec.artifacts` lists all the components which the profile will install.

The following artifact types are supported:
- ['Local' Helm Chart](/docs/author-docs/local-helm-chart)
- ['Remote' Helm Chart](/docs/author-docs/remote-helm-chart)
- [Raw Kubernetes yaml](/docs/author-docs/kustomize-and-raw)
- [Kustomize patch](/docs/author-docs/kustomize-and-raw)
- [Another profile](/docs/author-docs/nested-profiles)

Please refer to their dedicated docs pages for details on how to register different artifact
types in a profile.

:::info
Please take care to name your profile artifacts sensibly e.g. `fluent-bit`, `nginx-server`. 
This will help those who are installing your profile locate and refer to artifacts later.

See [how users use artifact names to configure values](/docs/installer-docs/setting-values)
for an example of how artifact names are used.
:::

## Profile repo directories

It will be assumed that everything contained within the same directory as a `profile.yaml`
is related to that same profile.

A repository can contain multiple profiles when they are written in separate directories.
For example, the following structure shows a repo with three distinct profiles which
can be installed independently of each other:

```bash
org-profiles-repo/
├── logging-profile
│   └── profile.yaml
├── observability-profile
│   └── profile.yaml
└── our-awesome-apps
    └── profile.yaml
```

:::tip
The name of each profile directory **must** match the name given in the `profile.yaml`
`metadata.name`.
:::

A repository can also contain just a single profile, with the `profile.yaml`
defined at the top level:

```bash
org-profiles-repo/
└── profile.yaml
```

Profile directories can contain other objects related to various artifacts. These
will be demonstrated in subsequent pages.

Examples of profiles with various artifacts and configurations can be found [here](https://github.com/weaveworks/profiles-examples).

## Documenting profiles

When you add ['local' Helm](/docs/author-docs/local-helm-chart), [raw yaml or Kustomize](/docs/author-docs/kustomize-and-raw)
artifacts to your profile, we recommend noting that you have done
so in your Readme, or other documentation. Users of such profiles will have to provide additional flags
when installing.

Take care to also note whether you are adding a [nested profile](/docs/author-docs/nested-profiles) which contains local resources
such as the ones listed above.

Users are also able to configure values on Helm artifacts. To help them discover which values
are available, you can provide information or links to the Charts you have used in your profile.
