---
sidebar_position: 3
---

# Install a profile

:::info

This stage of the tutorial assumes you have prepared your environment correctly.

Please refer back to the [set up docs](/docs/tutorial-basics/setup) if not.
:::

To install a profile, we use `pctl add`.

With the following command `pctl add` will:
- generate a set of manifests for each profile artifact
- commit those manifests to a branch in your GitOps repo
- push that branch and
- open a PR to merge the changes

Your GitOps repository is the one you synced to Flux in your cluster in the
[environment setup](/docs/tutorial-basics/setup#a-github-repo-synced-to-flux) section of this tutorial.

_(A breakdown of each flag is provided below.)_

```bash
pctl add \
  --name <profile installation name> \
  --profile-repo-url <URL of repo containing profile to add> \
  --profile-path . \
  --create-pr \
  --pr-repo <gitops repo username or orgname>/<gitops repo name> \
  --pr-branch add-simple-profile
```

Above we use the following flags:
- `--name`. This is the name of the profile installation. The installation directory will be created using this name.
- `--profile-repo-url`. This is the full URL of the repository containing the profile you wish to install on your cluster.
  If you completed the previous section and wrote your own profile, you can use that here.
  If you chose not to, you can use the following URL to a repository containing a simple example profile which
  will install an NGINX server: https://github.com/weaveworks/nginx-profile
- `--profile-path`. This is the relative path within the profile definition repo which contains the
  `profile.yaml`. Upstream profile repos can contain multiple profiles separated into
  different subdirectories. In our case, whether you are using the example repo, or the
  one you created in the previous section, there is just one profile, and the `profile.yaml`
  is located at the top level: `.`.
- `--create-pr`. This directs `pctl` to open a PR against the main branch of your GitOps repo.
  _Note that this flag is only supported for GitHub._
- `--pr-repo`. The partial URL of your GitOps repo synced to your cluster, in the format
  `username/repo-name`.
- `--pr-branch`. The name of the branch `pctl` will create in your GitOps repo to push
  changes to and open a PR against your main branch.

:::caution Private repos
If either your GitOps repository, or the repository containing the profile you wish to install
are private, remember to ensure that your local git environment is configured correctly.
:::

Once you have run the command, navigate to your GitOps repo and approve the PR.
Flux will then sync the new files, and the profile will be applied to your cluster.

You can eventually see the profile artifact running in the `default` namespace of your cluster
by running `kubectl get pod/nginx-server`.

To delete the profile, there is no need to use `pctl`. Simply remove the generated files from
your GitOps repo, merge the changes, and wait for Flux to delete those resources.

----------------------------------

:sparkles: The "Getting Started" tutorial ends here, please consult the documentation for more
advanced usage of profiles and `pctl`. :sparkles:
