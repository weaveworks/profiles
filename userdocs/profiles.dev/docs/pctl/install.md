---
sidebar_position: 2
---

# Install

<!-- TODO autogen-->

```sh
NAME:
   pctl install - generate a profile installation

USAGE:
   To install from a profile catalog entry: pctl --catalog-url <URL> install --name pctl-profile --namespace default --profile-branch main --config-map configmap-name <CATALOG>/<PROFILE>[/<VERSION>]
   To install directly from a profile repository: pctl install --name pctl-profile --namespace default --profile-branch development --profile-url https://github.com/weaveworks/profiles-examples --profile-path bitnami-nginx

OPTIONS:
   --name value            The name of the installation. (default: pctl-profile)
   --namespace value       The namespace to use for generating resources. (default: default)
   --profile-branch value  The branch to use on the repository in which the profile is. (default: main)
   --config-map value      The name of the ConfigMap which contains values for this profile.
   --create-pr             If given, install will create a PR for the modifications it outputs. (default: false)
   --pr-remote value       The remote to push the branch to. (default: origin)
   --pr-base value         The base branch to open a PR against. (default: main)
   --pr-branch value       The branch to create the PR from. Generated if not set.
   --out value             Optional location to create the profile installation folder in. This should be relative to the current working directory. (default: current)
   --pr-repo value         The repository to open a pr against. Format is: org/repo-name.
   --profile-url value     Optional value defining the URL of the profile.
   --profile-path value    Value defining the path to a profile when url is provided. (default: <root>)
   --git-repository value  The namespace and name of the GitRepository object governing the flux repo.
   --help, -h              show help (default: false)
```
