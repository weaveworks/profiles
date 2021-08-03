---
sidebar_position: 3
---

# Local Helm Chart artifacts

A "local" Helm Chart artifact means the chart itself is stored locally in the profile directory.
On installation these local manifests will be processed for addition to the GitOps repo, rather
than fetched down from a remote chart server.

Profile Authors may wish to add remote Charts to their profile repos if they:
- Are installing on a cluster with no internet access
- Are installing from a private Helm repository
- Wish to modify the chart themselves
- Would simply like to vendor their dependencies

_(If you do not believe you require a "local" Helm Chart artifact, please refer to
the page on [remote Helm Chart](/docs/author-docs/remote-helm-chart) artifacts.)_

:::tip
When you add local artifacts (meaning those with manifests stored in the profile repository)
to your profile, we recommend noting that you have done so in your Readme, or other documentation.
Users of such profiles will have to provide additional flags when installing.

Take care to also note whether you are adding a nested profile which contains local resources.
:::

To use local Helm resources, store the chart in a subdirectory within the profile
directory. Taking a previous structure example:

```bash
org-profiles-repo/
├── logging-profile
│   ├── fluentd
│   │   └── chart
│   │       ├── Chart.yaml
│   │       └── ...
│   └── profile.yaml
...
```

In the `profile.yaml` we add this local artifact by using the `path` key under the `chart`
type identifier:

```yaml
# ...
spec:
  # ...
  artifacts:
    - name: # the name of your artifact as you would like it to be known in the profile
      chart:
        path: "fluentd/chart" # the relative path to the chart directory
	# ...
```

The exact directory structure can be as you wish, as long as it is a child to the profile
directory and the `path` value in the `profile.yaml` is correct.

:::info
When `path` is present, any other fields under `chart` are ignored.
:::

If you would like your Chart artifact to be installed with some default values applied,
see the page [here](/docs/author-docs/default-values).

Examples of profiles with various artifacts and configurations can be found [here](https://github.com/weaveworks/profiles-examples).
