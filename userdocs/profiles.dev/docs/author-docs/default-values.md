---
sidebar_position: 6
---

# Default artifact values

Authors have the option to set default values on an artifact which will be applied
when users install their profile.

## Setting default values for Helm Chart artifacts

If you would like default configuration values to be applied to a Helm Chart artifact
when it is installed by a user, you can set the `defaultValues` field:

```yaml
  # ...
  artifacts:
    - name: foobar
      chart:
	    # ...
        defaultValues: |
          replicaCount: 3
          service:
            type: ClusterIP
	# ...
```

The value type for `defaultValues` is a `string`. When installed by a user, this data
will be placed directly into a ConfigMap and then applied to the user's cluster.

:::info
These values can be overridden by a user.
:::

Currently it is only possible to configure individual Chart artifacts. Soon authors
will be able to set "global" default variables which can be applied to multiple
or all Chart artifacts.

## Setting default values for other artifacts

This functionality is not yet available, but we are working on it!
