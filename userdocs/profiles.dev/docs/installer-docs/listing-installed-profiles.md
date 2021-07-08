---
sidebar_position: 5
---

# Listing installed profiles

You can see which profiles are installed on your cluster with the `pctl list` subcommand.

```bash
$ pctl list
NAMESPACE       NAME            SOURCE                                                                          AVAILABLE UPDATES
default         pctl-profile    nginx-catalog/weaveworks-nginx/v0.1.0                                           -
default         update-profile  https://github.com/weaveworks/profiles-examples:branch-and-url:bitnami-nginx    -
```

_In case of a branch install, as seen on the second line above, the source is put together as follows: `url:branch:profile-name`._

If you have installed your profiles via a catalog, you will be able to see whether updates are available:

```bash
$ pctl list
NAMESPACE       NAME            SOURCE                                  AVAILABLE UPDATES
default         pctl-profile    nginx-catalog/weaveworks-nginx/v0.1.0   v0.1.1
```

_There is currently no native way to update profiles, but we are working on it!_
