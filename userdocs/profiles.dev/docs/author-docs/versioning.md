---
sidebar_position: 9
---

# Versioning

Profiles are versioned via tags.

The pattern followed is one which already exists in the Kubernetes ecosystem,
for example in [Kustomize](https://github.com/kubernetes-sigs/kustomize/tags).

Tags can be specified in the following format: `<profile-name>/<version>`.
This allows profile authors to create multiple profiles in just one repository.

:::tip
Each tag must be valid semver.
:::

For example, in a profile repository containing the following profiles:

```bash
org-profiles-repo/
├── logging-profile
│   └── profile.yaml
├── observability-profile
│   └── profile.yaml
└── our-awesome-apps
    └── profile.yaml
```

We would create a new release version of the `logging-profile` by creating the tag
`logging-profile/v0.0.1`. The `observability-profile` would be versioned as `observability-profile/v0.0.1`,
and lastly `our-awesome-apps` would be `our-awesome-apps/v0.0.1`.

:::tip
Our tagging only supports going one level deep.

`profile-name/1.0.0` would therefore be valid, and `another-level/profile-name-1/0.0.1` would not.
:::

Authors can still create repositories which contain a single profile at the top level.
Tagging is completely optional, and users have methods of using your profile which do
not require tags.

Examples of tagged profiles can be found [here](https://github.com/weaveworks/profiles-examples).
