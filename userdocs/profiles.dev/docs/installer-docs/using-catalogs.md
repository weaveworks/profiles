---
sidebar_position: 3
---

# Using catalogs

If your cluster admin has added a profile catalog to your cluster, you can
use `pctl` to search for and install profiles approved by your organisation.

## Searching the catalog

To search the catalog for available profiles which match a query, call:

```bash
pctl search <query>
```

For example, to search for all profiles which would install an NGINX server:

```bash
$ pctl search nginx
CATALOG/PROFILE                 VERSION DESCRIPTION
nginx-catalog/bitnami-nginx     v0.0.2  Profile for deploying local nginx chart
nginx-catalog/weaveworks-nginx  v0.1.0  Profile for deploying nginx
...
```

To see all available profiles, just pass the `--all` flag:

```bash
pctl search --all
```

## Inspecting profiles in the catalog

To learn more about a particular profile, use the `show` subcommand:

```bash
$ pctl show nginx-catalog/bitnami-nginx
Catalog         nginx-catalog
Name            bitnami-nginx
Version         v0.0.2
Description     Profile for deploying local nginx chart
URL             https://github.com/weaveworks/profiles-examples
Maintainer      weaveworks
Prerequisites   kubernetes 1.19
```

_Note that the Prerequisites field is not yet processed, we are working on it!_

## Installing a profile from the catalog

To install a profile from the catalog we provide a positional argument after all other flags
in the format of `<catalog name>/<profile name>`.

```bash
pctl add \
  --create-pr \
  --pr-repo <gitops repo username or orgname>/<gitops repo name> \
  nginx-catalog/weaveworks-nginx
```

The above command will install the latest version.
To install a specific version of a profile, simply add it to the end:

```bash
pctl add \
  --create-pr \
  --pr-repo <gitops repo username or orgname>/<gitops repo name> \
  nginx-catalog/weaveworks-nginx/v0.0.2
```

:::info
We recommended also setting the `--git-repository` flag. See [the section here](/docs/installer-docs/installing-via-gitops#the-git-repository-flag)
for more information.
:::
