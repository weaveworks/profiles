---
sidebar_position: 6
---

# Removing profiles

To uninstall any profile from your cluster, just remove the directory from your
GitOps repo, and open a PR with that commit.

After you approve and merge the PR, Flux will sync the changes and the profile will
no longer be running in your cluster.
