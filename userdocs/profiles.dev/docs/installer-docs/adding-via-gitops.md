---
sidebar_position: 1
---

# Adding Profile with GitOps

## Environment setup

To add profiles with GitOps you will need to have completed the following environment setup.

 :white_check_mark: &nbsp;&nbsp; [Cluster](/docs/tutorial-basics/setup#kubernetes-cluster)

 :white_check_mark: &nbsp;&nbsp; [Pctl binary](/docs/tutorial-basics/setup#pctl-the-profiles-cli)

 :white_check_mark: &nbsp;&nbsp; [Profiles CRDs and Flux CRDs](/docs/tutorial-basics/setup#profiles-crds-and-flux-crds)

 :white_check_mark: &nbsp;&nbsp; [GitHub repo](/docs/tutorial-basics/setup#a-github-repo-synced-to-flux)

 :white_check_mark: &nbsp;&nbsp; [GitHub token](/docs/tutorial-basics/setup#personal-access-token)

The full setup docs can be found [here](/docs/tutorial-basics/setup#prerequisites).

## Simple add from a profile URL

:::caution Private repos
If either your GitOps repository, or the repository containing the profile you wish to install
are private, remember to ensure that your local git environment is configured correctly.
:::

To add a profile, we use `pctl add`. To see all flags available on this subcommand,
see [the help](/docs/pctl/pctl-add-cmd).

There are two methods by which you can add a profile: with a direct URL and via a catalog.
Please see [the page here](/docs/installer-docs/using-catalogs) for specific instructions on how to add a profile from a catalog.

With the following command, `pctl` will:
- generate a set of manifests for each artifact declared in the profile at the given URL
- commit those manifests to a branch in your GitOps repo
- push that branch and
- open a PR on your GitOps repo to merge the changes.
Your GitOps repo is the one you synced to Flux in your cluster in the
[environment setup](/docs/tutorial-basics/setup#a-github-repo-synced-to-flux).

```bash
pctl add \
  --profile-url <URL of profile to install> \
  --create-pr \
  --pr-repo <gitops repo username or orgname>/<gitops repo name>
```

Above we use the following flags:
- `--profile-url`. This is the full git URL of the profile you wish to install on your cluster.
- `--create-pr`. This directs `pctl` to open a PR against the main branch of your GitOps repo.
  _Note that this flag is only supported for GitHub._
- `--pr-repo`. The partial URL of your GitOps repo synced to your cluster, in the format
  `username/repo-name` (or `org-name/repo-name`).

:::info
We recommended also setting the `--git-repository` flag. See [below](#the-git-repository-flag)
for more information.
:::

Once you have run the command, navigate to your GitOps repo and approve and merge the PR.
Flux will then sync the new files, and the profile will be applied to your cluster.

:::info
When installing a profile via its URL (i.e. when using the `--profile-url` flag)
remember to check where the profile's `profile.yaml` file is located within
the profile's source repository.

Once discovered, you can set the relative path to this file using the `--profile-path` flag.
:::

## Further configurations

You can pass further arguments to `pctl add` for more control over the format
and destination of your PR. For example, the following command will:
- generate the profile manifests in the `ethel` directory within my GitOps repo
- set the objects to be deployed in the `two-sugars-please` namespace in my cluster
- commit to a branch called `add-aunt-ethels-profile`
- push that branch to the `jammy-dodger` remote
- open a PR against the `tea-time` branch in my target GitOps repo

```bash
pctl add \
  --profile-url https://github.com/weaveworks/nginx-profile \
  --create-pr \
  --pr-repo drwho/thirteen \
  --out ethel \
  --namespace two-sugars-please \
  --pr-branch add-aunt-ethels-profile \
  --pr-remote jammy-dodger \
  --pr-base tea-time
```

To see all flags available on this subcommand, see [the help](/docs/pctl/pctl-add-cmd).

## The `git-repository` flag

:::tip
While not required, it is recommended to always set the `--git-repository` flag.
This flag is needed by the majority of artifact types, so it may be easier to always
set it rather than have to discover whether a profile requires it.

Keep reading for instructions on how to set this option.
:::

This flag is required for installing profiles with local resources.
A local resource is any profile artifact which refers to manifests stored within the profile
repository. You can see [this page](/docs/author-docs/local-helm-chart) from the author docs as an example.

In order to install profiles with local resources, you must provide `pctl` with information
about your GitOps repo (very meta, I know).

First, look up the ID of the `GitRepository` resource connected to your repo. (This is
what Flux uses internally to keep things up to date between the repo and the cluster.)

```bash
# replace GITOPS_REPO with the name of your GitOps repository

kubectl get gitrepositories.source.toolkit.fluxcd.io -A | awk '/GITOPS_REPO/ {print $1"/"$2}'
```

This should return something like:

```bash
# namespace/name
flux-system/flux-system
```

Then add the following flag to your `add` command:

```sh
--git-repository <output from the command above>
```

### When to use the `--git-repository` flag

As a general rule, we recommended always setting this flag, but if you're curious to find out
you can navigate to the repository of the profile and look at the `profile.yaml`.
In this file you will find a list of `artifacts`. Those which have local resources can be identified
with one of the following keys:
- `chart.path`
- `kustomize.path`

If the list of artifacts includes any `profile.source` keys,
these may point to other profiles which themselves could contain local resources.

Profile authors may also note in their documentation that their profile contains
local resources, but this should not be relied upon.
