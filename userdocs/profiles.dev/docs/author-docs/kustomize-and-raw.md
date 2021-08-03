---
sidebar_position: 4
---

# Kustomize patches and raw YAML artifacts

Both Kustomize patches and raw YAMLs (such as a simple deployment manifest,
or any other Kubernetes object) can be added to a profile under the same key.

The resources must be stored locally within the profile directory.
For example:

```bash
org-profiles-repo/
├── our-awesome-apps-profile
│   ├── super-cool-artifact-manifests
│   │   ├── deployment.yaml
│   │   └── patches.yaml
│   └── profile.yaml
...
```

Then in the `profile.yaml` we add these artifact manifests by using the `kustomize`
type identifier:

```yaml
# ...
spec:
  # ...
  artifacts:
    - name: # the name of your artifact as you would like it to be known in the profile
      kustomize:
        path: "super-cool-artifact-manifests/" # the relative path to the manifests directory
    # ...
```

:::tip
When you add local artifacts (meaning those with manifests stored in the profile repository)
to your profile, we recommend noting that you have done so in your Readme, or other documentation.
Users of such profiles will have to provide additional flags when installing.

Take care to also note whether you are adding a nested profile which contains local resources.
:::

The exact directory structure can be as you wish, as long as it is a child to the profile
directory and the `kustomize.path` value in the `profile.yaml` is correct.

Examples of profiles with various artifacts and configurations can be found [here](https://github.com/weaveworks/profiles-examples).
