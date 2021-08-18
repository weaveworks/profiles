---
sidebar_position: 1
---

# Installing Profile with GitOps

## Environment setup

To install profiles with GitOps you will need to have completed the following environment setup.

 :white_check_mark: &nbsp;&nbsp; [Cluster](/docs/tutorial-basics/setup#kubernetes-cluster)

 :white_check_mark: &nbsp;&nbsp; [Pctl binary](/docs/tutorial-basics/setup#pctl-the-profiles-cli)

 :white_check_mark: &nbsp;&nbsp; [Profiles CRDs and Flux CRDs](/docs/tutorial-basics/setup#profiles-crds-and-flux-crds)

 :white_check_mark: &nbsp;&nbsp; [GitHub repo](/docs/tutorial-basics/setup#a-github-repo-synced-to-flux)

 :white_check_mark: &nbsp;&nbsp; [GitHub token](/docs/tutorial-basics/setup#personal-access-token)

The full setup docs can be found [here](/docs/tutorial-basics/setup#prerequisites).

## Simple install from a profile URL


### Bootstrapping your local git repository
You can use the `pctl bootstrap` command to save commonly used `pctl` configuration to your GitOps repo.

Once such piece of configuration is your flux repo's `GitRepository` resource.
If you have **not** bootstrapped your local GitOps repository, you will have to provide the `--git-repository` flag when installing profiles (see [below](#the-git-repository-flag) for more detail). The `--git-repository` references the namespace and name of the
[Flux `GitRepository`](https://fluxcd.io/docs/components/source/gitrepositories/)
resource that is pointing at your GitOps repository. The value should be in the format `<namespace>/<name>`, for example
`flux-system/gitops-repo`. This value is needed in order for pctl to generate Flux resources, such as `Kustomization`s.
To bootstrap your local git repository run the following:

```bash
pctl bootstrap --git-repository flux-system/gitops-repo ~/workspace/gitops-repo/
```

_See [below](#the-git-repository-flag) for how to find the namespace and name of your `GitRepository` resource._

 This will create a `.pctl/config.yaml` in your git repository to store the value for `--git-repository`. Future
`pctl add` commands in this repository will then detect the pctl config and re-use the value.

### Installing

:::caution Private repos
If either your GitOps repository, or the repository containing the profile you wish to install
are private, remember to ensure that your local git environment is configured correctly.
:::

To install a profile, we use `pctl add`. To see all flags available on this subcommand,
see [the help](/docs/pctl/pctl-add-cmd).

There are two methods by which you can install a profile: with a direct URL and via a catalog.
Please see [the page here](/docs/installer-docs/using-catalogs) for specific instructions on how to
install a profile from a catalog.

With the following command, `pctl` will:
- generate a set of manifests for each artifact declared in the profile at the given URL
- commit those manifests to a branch in your GitOps repo
- push that branch and
- open a PR on your GitOps repo to merge the changes.
Your GitOps repo is the one you synced to Flux in your cluster in the
[environment setup](/docs/tutorial-basics/setup#a-github-repo-synced-to-flux).

```bash
pctl add \
  --name <name of the profile installation> \
  --profile-repo-url <URL of repo containing profile to install> \
  --create-pr \
  --pr-repo <gitops repo username or orgname>/<gitops repo name>
```

Above we use the following flags:
- `--name`. This is the name of the profile installation. The installation directory will be created using this name.
- `--profile-repo-url`. This is the full URL of the repository containing the profile you wish to install on your cluster.
- `--create-pr`. This directs `pctl` to open a PR against the main branch of your GitOps repo.
  _Note that this flag is only supported for GitHub._
- `--pr-repo`. The partial URL of your GitOps repo synced to your cluster, in the format
  `username/repo-name` (or `org-name/repo-name`).

:::info
If you have not bootstrapped your local git repository you must provide the`--git-repository` flag.
See [below](#the-git-repository-flag)
for more information.
:::

Once you have run the command, navigate to your GitOps repo and approve and merge the PR.
Flux will then sync the new files, and the profile will be applied to your cluster.

:::info
When installing a profile via its repository's URL (i.e. when using the `--profile-repo-url` flag)
remember to check where the profile's `profile.yaml` file is located within
the profile's source repository.

Once discovered, you can set the relative path to this file using the `--profile-path` flag, which defaults to the root of the repository.
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
  --name pctl-profile \
  --profile-repo-url https://github.com/weaveworks/nginx-profile \
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

This flag is required for installing profiles because pctl generates flux `Kustomization` resources for
deploying the profile artifacts, and these resources need to know which `GitRepository` resource is governing the repository.
To avoid setting this flag on every call of `pctl add`, users can bootstrap the repository; see [here](#bootstrapping-your-local-git-repository).
Any `pctl add` call from within this repository will then no longer need this flag.

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
