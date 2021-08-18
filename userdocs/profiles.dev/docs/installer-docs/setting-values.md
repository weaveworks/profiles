---
sidebar_position: 2
---

# Configuring values

Users can partially configure their profile installations with the `--config-map`
flag on the `pctl add` command.

## Helm Chart artifact values

To configure Helm Chart artifacts in a profile, first create a [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/).

For example, to configure a profile with the following artifacts:

```yaml
# ...
spec:
  # ...
  artifacts:
  - name: fluent-bit
    chart:
      path: "fluent-bit/chart"
  - name: nginx-server
    chart:
      url: https://charts.bitnami.com/bitnami
      name: nginx
      version: "8.9.1"
  # ...
```

You can create a ConfigMap like so:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-profile-values
  namespace: default
data:
  fluent-bit: |
    replicaCount: 3
  nginx-server: |
    service:
      type: ClusterIP
```

:::tip
Note that the data fields in the ConfigMap **must** match the names of the artifacts
you are configuring.
:::

Commit your ConfigMap yaml to your GitOps repository so that Flux can sync it to your cluster.
You can then provide the name of your ConfigMap to the `add` command:

```bash
pctl add \
  --name <profile installation name> \
  --profile-repo-url <URL of repo containing profile to install> \
  --create-pr \
  --pr-repo <gitops repo username or orgname>/<gitops repo name> \
  --config-map my-profile-values
```

When you have merged the PR that `pctl` opens in your GitOps repo, your profile will be deployed
with your settings applied.

:::info
To discover which values you are able to set on Helm Chart Artifacts, you will have to look up
each Chart artifact listed in the profile.

We encourage profile Authors to provide details or links to configurable values,
so please check any profile documentation.
:::

## Configuring multiple Helm artifacts with one value

This is not possible with the current version of profiles/pctl, but we are working
on a way for users to specify values required by all artifacts (e.g. `hostname` or `clustername`)
across multiple Helm Chart artifacts. Watch this space!

## Other profile artifacts

There is currently no way to configure values of other artifacts, but we are working on it!
