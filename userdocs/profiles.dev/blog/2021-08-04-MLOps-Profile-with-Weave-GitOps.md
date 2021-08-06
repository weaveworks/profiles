---
slug: mlops-profile
title: Install MLOps Profile with Weave GitOps
author: Chanwit Kaewkasi
author_title: DX Engineer, WeaveWorks
tags: [mlops, profile, gitops]
---

“_Profiles_” is a GitOps-native package management system, which allows you to apply a Profile to add a set of capabilities to your GitOps-enabled cluster.  Here’s a quick tutorial on how to provision a cluster, install Weave GitOps and apply the MLOps profile to enable Kubeflow Pipeline for your cluster. To try this setup, we strongly recommend that you use [Kind](https://github.com/kubernetes-sigs/kind) as a local cluster on your laptop.

<!--truncate-->

First, you need Weave GitOps Core. Please download it from [here](https://docs.gitops.weave.works/docs/installation/)

```
curl -L "https://github.com/weaveworks/weave-gitops/releases/download/v0.2.2/wego-$(uname)-$(uname -m)" -o wego
chmod +x wego
sudo mv ./wego /usr/local/bin/wego
wego version
```

Second, you also need Profiles CLI. You can download the CLI from [here](https://github.com/weaveworks/pctl/releases)

The following are links of the Profile CLI for both Linux and macOS.

- [Profiles CLI for Linux](https://github.com/weaveworks/pctl/releases/download/v0.2.0/pctl_linux_amd64.tar.gz)
- [Profiles CLI for macOS](https://github.com/weaveworks/pctl/releases/download/v0.2.0/pctl_darwin_amd64.tar.gz)

#### Here’s steps to install MLOps Profile with Weave GitOps

1. Create cluster

```
kind create cluster
```

2. Create a new GitOps repository on GitHub.

For example, my GitOps repo is [test-mlops-profile](https://github.com/chanwit/test-mlops-profile)

Obtain your GitHub token, and export it as an environment variable along with your user and repo names. Please change `GITHUB_USER` to your GitHub user name.

```
export GITHUB_TOKEN=gh0_1234567890
export GITHUB_USER=chanwit
export GITHUB_REPO=test-mlops-profile
```

3. Check that you already installed both Weave GitOps Core and Profiles CLI, then we install GitOps runtime to our cluster using the following command.
```
wego gitops install
```

5. After that, we install the ProfileManager controller using the pctl install command. 
Please note that we need to tell it that our GitOps runtime is running in the wego-system namespace.

```
pctl install --flux-namespace=wego-system
```

6. Then, we clone the GitOps repo and setup our Weave GitOps automation for it.
```
git clone git@github.com:$GITHUB_USER/$GITHUB_REPO
cd $GITHUB_REPO
wego app add .
```

7. Weave GitOps will create a PR for you to merge. After you merge it, GitOps automation will kick off.

8. Finally, we are able to install the MLOps profile using the following commands.
```
kubectl create ns kubeflow

pctl add --profile-url=https://github.com/chanwit/mlops-profile --git-repository=wego-system/test-mlops-profile --name=mlops-profile --namespace=kubeflow

git add .
git commit -am “install mlops profile”
git push origin main
```

Here’s the pod list after you have the MLOps profile installed.

![MLOps profile](/img/mlops.png)
