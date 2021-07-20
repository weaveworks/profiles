---
sidebar_position: 9
---

# Dependencies between artifacts

It's possible to add dependencies between artifacts. Consider the following scenario:

There is an artifact (`artifact B`) which needs a certain set of CRDs to be present on the cluster before it can be installed.
In that case, `artifact B` can depend on `artifact A` which installs the necessary CRDs before `artifact B`.

This is achieved by adding a `dependsOn` setting on the desired artifact's configuration.

Such as:

```yaml
spec:
  description: Depending artifacts
  artifacts:
    - name: artifact-a
      kustomize:
        path: crds/rbacs
    - name: artifact-b
      chart:
        url: https://charts.bitnami.com/bitnami
        name: nginx
        version: "9.3.0"
        defaultValues: |
          service:
            type: ClusterIP
      dependsOn:
        - name: artifact-a
        # ... the name of any other further dependencies this artifact might have 
```

:::info
At the time of this writing, nested-profile dependencies are not functional yet.
:::

_Note_: this uses flux's [dependsOn](https://fluxcd.io/docs/components/kustomize/kustomization/#kustomization-dependencies) mechanism in the background to make sure artifacts
are created in order. What this brings with it is, that every Kubernetes object is wrapped into a `Kustomization` object such as:

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1beta1
kind: Kustomization
metadata:
  name: pctl-profile-weaveworks-nginx-dependon-chart
  namespace: default
spec:
  dependsOn:
  - name: pctl-profile-weaveworks-nginx-nginx-deployment
    namespace: default
  - name: pctl-profile-weaveworks-nginx-nginx-chart
    namespace: default
  healthChecks:
  - apiVersion: helm.toolkit.fluxcd.io/v2beta1
    kind: HelmRelease
    name: pctl-profile-weaveworks-nginx-dependon-chart
    namespace: default
  interval: 5m0s
  path: artifacts/dependon-chart/helm-chart
  prune: true
  sourceRef:
    kind: GitRepository
    name: flux-system
    namespace: flux-system
  targetNamespace: default
```

Notice the `dependsOn` section. This is how flux controls this resource's dependency chain. The actual HelmRelease is handled
by this `Kustomization` resource via the `healthChecks` attribute.

This also means that every profile that is installed will have two files in the installation folder.

- kustomization.yaml
- kustomize-flux.yaml

The kustomization.yaml file contains this:
```yaml
resources:
- kustomize-flux.yaml
```

This, tells flux to only care about the resource that is inside the kustomize-flux.yaml file.

Examples of profiles with various artifacts and configurations can be found [here](https://github.com/weaveworks/profiles-examples).
