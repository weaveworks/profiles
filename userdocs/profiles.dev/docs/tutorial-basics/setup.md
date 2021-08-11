---
sidebar_position: 1
---

# Environment setup

:::info Assumed knowledge

This tutorial assumes you have some knowledge of the concept of GitOps and are comfortable using
[Flux](https://fluxcd.io/).

Please refer to the [Introduction](/docs/intro) to read about the core concepts of Profiles.
:::



In this tutorial you will create and install a simple profile onto your Kubernetes cluster using various GitOps tools.

_If you are only interested in **installing** profiles, not writing them, please skip ahead to the relevant section
once you have set up your environment._

_If you are only interested in **writing** profiles, not installing them, you can skip the environment
setup steps._

------------------

## Prerequisites

In order to install profiles, you need to have the following set up:

### Kubernetes cluster

For local testing, we recommend using [kind](https://kind.sigs.k8s.io/docs/user/quick-start/).
The cluster must be version 1.16 or newer.

### Profiles CLI

Profiles are installed and managed via the official CLI: `pctl`.
Releases can be found [here](https://github.com/weaveworks/pctl/releases).
`pctl` binaries are not backwards compatible, and we recommended keeping your local
version regularly updated.

### Profiles CRDs and Flux CRDs

Profiles relies on Flux to deploy artifacts to your cluster, this means that at a minimum
you much have the following Flux CRDs and associated controllers installed:

- `helmreleases.helm.toolkit.fluxcd.io`
- `gitrepositories.source.toolkit.fluxcd.io`
- `helmrepositories.source.toolkit.fluxcd.io`
- `kustomizations.kustomize.toolkit.fluxcd.io`

You can install everything by running Flux's [install command](https://fluxcd.io/docs/cmd/flux_install/):

```bash
flux install
```

Or to install specific components:

```bash
flux install --components="source-controller,kustomize-controller,helm-controller"
```

Next install the Profiles CRD, with:

```bash
pctl install
```

Note: This will install the latest version of the profiles CRD, which may not always be stable.

To specify a [specific version](https://github.com/weaveworks/profiles/releases), use the `--version` flag.

### A GitHub repo, synced to Flux

This tutorial will require a GitHub account. (More git providers will be added in the future.)

The repo can be public or private (note: you will not be asked to push any sensitive information) and must
be linked to the Flux instance running in your cluster.

You can do this by running [`flux bootstrap github`](https://fluxcd.io/docs/installation/#github-and-github-enterprise) with the relevant arguments.

:::caution Private repos
If you choose to use a private repo, please ensure that your local git environment is set
up correctly for the rest of the tutorial.
:::

### Personal Access Token

The profile will be installed in a GitOps way, therefore `pctl` will push all manifests to your cluster git repo.
Create a [personal access token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line) (check all permissions under repo)
on your GitHub account and export it:

```bash
export GITHUB_TOKEN=<your token>
```

## Get started!

Check you have everything on this list and go back if something is missing.

 :white_check_mark: &nbsp;&nbsp; [Cluster](#kubernetes-cluster)

 :white_check_mark: &nbsp;&nbsp; [Pctl binary](#pctl-the-profiles-cli)

 :white_check_mark: &nbsp;&nbsp; [Profiles CRDs and Flux CRDs](#profiles-crds-and-flux-crds)

 :white_check_mark: &nbsp;&nbsp; [GitHub repo](#a-github-repo-synced-to-flux)

 :white_check_mark: &nbsp;&nbsp; [GitHub token](#personal-access-token)

Once you have completed the prerequisites installation you can start writing a profile!
