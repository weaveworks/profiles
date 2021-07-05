---
sidebar_position: 2
---

# Remote Helm Chart artifacts

In the [Getting started tutorial](/docs/tutorial-basics/create-a-profile),
we wrote a simple profile which contained a single artifact.
This artifact was a 'remote' Helm Chart. By remote we mean that we did not
store the Chart, and any other relevant manifests, in the profile directory.

_(For details on how to add a "local" Helm Chart artifact to your profile,
see the page [here](/docs/author-docs/local-helm-chart).)_

To add a remote Helm Chart artifact, you add the following to the `artifacts` list
in your `profile.yaml` spec:

```yaml
# ...
spec:
  # ...
  artifacts:
    - name: # the name of your artifact as you would like it to be known in the profile
      chart:
        url: # the full URL to the chart repository server
        name: # the name of the chart
        version: # the version of the chart
	# ...
```

In the snippet above, the `chart` key denotes the type of the artifact you are
adding.

If you would like your Chart artifact to be installed with some default values applied,
see the page [here](/docs/author-docs/default-values).

Examples of profiles with various artifacts and configurations can be found [here](https://github.com/weaveworks/profiles-examples).
