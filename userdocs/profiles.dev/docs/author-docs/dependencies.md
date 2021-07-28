---
sidebar_position: 7
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
