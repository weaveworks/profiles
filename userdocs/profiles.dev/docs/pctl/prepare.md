---
sidebar_position: 1
---

# Prepare

<!-- TODO autogen-->

```sh
NAME:
   pctl prepare - prepares the cluster for profiles by deploying the profile controllers and custom resource definitions

USAGE:
   pctl prepare

OPTIONS:
   --dry-run                  If defined, nothing will be applied. (default: false)
   --keep                     Keep the downloaded manifest files. (default: false)
   --ignore-preflight-errors  Instead of stopping the process, output warnings when they occur during preflight check. (default: false)
   --version value            Define the tagged version to use which can be found under releases in the profiles repository. Exp: [v]0.0.1
   --baseurl value            Define the url to go and fetch releases from. (default: https://github.com/weaveworks/profiles/releases)
   --flux-namespace value     Define the namespace in which flux is installed. (default: flux-system)
   --out value                Specify the output location of the downloaded prepare file. (default: os.Temp)
   --context value            The Kubernetes context to use to apply the manifest files .
   --help, -h                 show help (default: false)
```
