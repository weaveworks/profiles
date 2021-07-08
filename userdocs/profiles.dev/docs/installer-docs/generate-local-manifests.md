---
sidebar_position: 4
---

# Generating local manifests

If you are curious to see what `pctl` will create _without_ opening a PR
in your GitOps repo, you can generate the files locally by dropping all the `pr`
related flags.

For example, using a profile URL:
```yaml
pctl install \
  --profile-url <URL of profile to install> \
  --out relative-path
```

:::info
When installing a profile via its URL (i.e. when using the `--profile-url` flag)
remember to check where the profile's `profile.yaml` file is located within
the profile's source repository.

Once discovered, you can set the relative path to this file using the `--profile-path` flag.
:::

Example generating from a profile listed in a catalog:

```yaml
pctl install \
  --out relative-path \
  <catalog name>/<profile>
```
