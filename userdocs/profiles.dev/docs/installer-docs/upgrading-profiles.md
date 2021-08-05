---
sidebar_position: 5
---

# Upgrading profiles

When newer versions of a profile are available you will be able to discover them by 
[listing installed profiles](/docs/installer-docs/listing-installed-profiles#listing-installed-profiles).
Once you know which version you want to upgrade to, run the following:

```bash
# The first argument is the relative or aboslute path to the local installation directory
# and the second argument is the version to upgrade to.
pctl upgrade ~/workspace/demo-profile/ v0.0.2
```

This will then perform an upgrade of your local installation. You can also pass in the `--create-pr` flag to automatically create a PR
. Pctl uses a 3-way merge behind the scenes to perform the upgrade. If you have made local modifications to
your installation that conflict with changes in the upgrades you will get merge conflicts, and will have to manually resolve them.
