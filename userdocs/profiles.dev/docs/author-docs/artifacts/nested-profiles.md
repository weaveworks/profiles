---
sidebar_position: 5
---

# Nested profile artifacts

Profiles can refer to other profiles in their list of artifacts.

Simply use the `profile` field:

```yaml
# ...
spec:
  # ...
  artifacts:
    - name: # the name of your artifact as you would like it to be known in the profile
      profile:
        source:
          url: # required: fully qualified URL to the nested profile repository
          branch: # optional: the repo branch if the profile is somewhere other than `main`
          path: # optional: the relative path to a profile's directory within the repo
          tag: # optional: the tag of the profile
    # ...
```

:::danger
Do not attempt to add a profile into its own list of artifacts.

Likewise, please avoid adding a nested profile which contains a reference
to this profile.

We will detect this recursion or "circular import" when someone tries to install the profile,
but try not to do it anyway :wink:
:::

Examples of profiles with various artifacts and configurations can be found [here](https://github.com/weaveworks/profiles-examples).
