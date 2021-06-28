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

_If you are not interested in writing profiles, just installing them, please skip ahead to the relevant section
once you have set up your environment._

------------------

## Prerequisites

In order to install profiles, you need to have the following set up:

### Kubernetes cluster

For local testing, we recommend using [kind](https://kind.sigs.k8s.io/docs/user/quick-start/).
The cluster must be version 1.16 or newer.

### Flux components

While profiles can be installed manually, it is recommended to install them in a GitOps fashion,
and that means using Flux.

The Flux binary is not required, but several CRDs must be present on the cluster.
These are:
- buckets.source.toolkit.fluxcd.io
- gitrepositories.source.toolkit.fluxcd.io
- helmcharts.source.toolkit.fluxcd.io
- helmreleases.helm.toolkit.fluxcd.io
- helmrepositories.source.toolkit.fluxcd.io
- kustomizations.kustomize.toolkit.fluxcd.io

For simplicity, we recommend that you install the standard set of Flux components into your cluster.
This can be done by running [`flux check --pre`](https://fluxcd.io/docs/cmd/flux_check/) followed by [`flux install`](https://fluxcd.io/docs/cmd/flux_install/).
Those with more familiarity with Flux are free to fine-tune their installation and only install
what is necessary.

### A GitHub repo, synced to Flux

This tutorial will require a GitHub account. (More git providers will be added soon.)

The repo can be public or private (note: you will not be asked to push any sensitive information) and must
be linked to the Flux instance running in your cluster.

You can do this by running [`flux bootstrap github`](https://fluxcd.io/docs/installation/#github-and-github-enterprise) with the relevant arguments.

### pctl: the Profiles CLI

Profiles are installed and managed via the official CLI: `pctl`.
Releases can be found [here](https://github.com/weaveworks/pctl/releases).
`pctl` binaries are not backwards compatible, and we recommended keeping your local
version regularly updated.

### Personal Access Token

The profile will be installed in a GitOps way, therefore `pctl` will push all manifests to your cluster git repo.
Create a [personal access token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line) (check all permissions under repo)
on your GitHub account and export it:

```bash
export GIT_TOKEN=<your token>
```

## Get started!

Check you have everything on this list and go back if something is missing.

 :white_check_mark: [Cluster](#kubernetes-cluster)

 :white_check_mark: [Flux CRDs](#flux-components)

 :white_check_mark: [GitHub repo](#a-github-repo-synced-to-flux)

 :white_check_mark: [Pctl binary](#pctl-the-profiles-cli)

 :white_check_mark: [GitHub token](#personal-access-token)

Once you have completed the prerequisites installation you can start writing a Profile!
