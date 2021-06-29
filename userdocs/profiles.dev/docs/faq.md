---
sidebar_position: 10
---

# FAQ

Q. **Can I use something other than GitHub to host my GitOps repo?**

A. At the moment, GitHub is the only git provider supported by `pctl`. BUT this only affects
   the `pctl install` command when using the `--create-pr` flag. If you wish to use, say, GitLab
   to manage your cluster, you can use `pctl` to generate the profile manifests to a local directory,
   and then commit them to your Flux-synced repo manually.
