---
sidebar_position: 4
---

# Generating local manifests

If you are curious to see what `pctl` will create _without_ opening a PR
in your GitOps repo, you can generate the files locally by dropping all the `pr`
related flags.

For example, using a profile URL:
```yaml
pctl add \
  --name <profile installation name> \
  --profile-repo-url <URL of repo containing profile to install> \
  --out <relative path>
```

:::info
When installing a profile via its repository's URL (i.e. when using the `--profile-repo-url` flag)
remember to check where the profile's `profile.yaml` file is located within
the profile's source repository.

Once discovered, you can set the relative path to this file using the `--profile-path` flag, which 
defaults to the root of the repository.
:::

Example generating from a profile listed in a catalog:

```yaml
pctl add \
  --name <profile installation name> \
  --out <relative path> \
  <catalog name>/<profile>
```

Consider the following installation:

```
pctl add --name nginx-profile --git-repository flux-system/flux-system nginx-catalog/nginx/v2.0.1
generating profile installation from source: catalog entry nginx-catalog/nginx/v2.0.1
installation completed successfully
```

Let's take a look inside:

```
tree nginx-profile
nginx-profile
├── artifacts
│   └── bitnami-nginx
│       ├── helm-chart
│       │   ├── ConfigMap.yaml
│       │   ├── HelmRelease.yaml
│       │   └── HelmRepository.yaml
│       ├── kustomization.yaml
│       └── kustomize-flux.yaml
└── profile-installation.yaml
```

This profile installs an NGINX load balancer using a remote helm chart.

The folders contain the following files in order:

* bitnami-nginx - the name of the artifact
* helm-chart - contains the resources which install the actual nginx
* ConfigMap.yaml - contains any default values set by the author on the artifact
* HelmRelease.yaml - contains the HelmRelease object which installs nginx
* HelmRepository.yaml - contains the definition of the helm chart repository where the chart is located
* kustomization.yaml - this is a file which tells flux what to install -- see below
* kustomize-flux.yaml - this is a Kustomization object which deals with [dependencies](/docs/author-docs/dependencies)
* profiles-installation.yaml - contains information about the profile -- mainly used by pctl

What is `kustomization.yaml`? This is to prevent flux installing whatever there is in the `helm-chart` folder. The `helm-chart`
folder can contain local resources such as helm chart definitions, READMEs and non-kubernetes objects. If flux were to try to
install those, it would fail.

`kustomization.yaml` contains a single resource line:

```yaml
resources:
- kustomize-flux.yaml
```

Which means flux will only install the resource defined in this file. The Kustomization object in `kustomize-flux.yaml`
will take care of installing HelmRelease.
