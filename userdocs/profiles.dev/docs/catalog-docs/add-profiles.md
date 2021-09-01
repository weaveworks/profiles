---
sidebar_position: 1
---

# Adding profiles

:::info
Please refer to the [Introduction](/docs/intro#catalog) for a list of terms
around Catalogs.
:::

To add profiles to a catalog, you must first deploy the catalog API
and `ProfileCatalogSource` controller to your cluster.

This can be done with the `pctl install` command. By default, this command
will install latest, but you can pass a specific version with the `--version`
flag.

Once deployed, there are two methods to add profiles to the cluster catalog:
dynamically and manually.

## Dynamically populate a catalog

To dynamically add all profiles and tags from profiles repositories to the catalog,
create a `ProfileCatalogSource` providing the repositories' URLs,
and apply the resource to the cluster.

For example, the following manifest will discover all the profiles within the
[weaveworks/profiles-examples](https://github.com/weaveworks/profiles-examples) repo
and add an entry in the catalog for each profile tag.

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileCatalogSource
metadata:
  name: nginx-catalog
spec:
  repositories:
  - url: https://github.com/weaveworks/profiles-examples
```

After applying the manifest and waiting a moment, we can use `pctl get --catalog` to see the
catalogued profiles:

```bash
$ kubectl apply -f dynamic-catalog-source.yaml
# allow a few moments. the more profiles/tags, the more time the catalog manager
# will need to discover them all

$ pctl get --catalog
CATALOG/PROFILE                 VERSION DESCRIPTION
nginx-catalog/bitnami-nginx     v0.0.2  Profile for deploying local nginx chart
nginx-catalog/weaveworks-nginx  v0.1.0  Profile for deploying nginx
nginx-catalog/bitnami-nginx     v0.0.1  Profile for deploying local nginx chart
nginx-catalog/weaveworks-nginx  v0.1.1  Profile for deploying nginx
```

Once added, the catalog will monitor each profile and update the catalog entries
when new versions are released.

### Adding profiles from private repositories

To dynamically add profiles from a private repository, you must provide a reference to a
secret in your catalog source manifest.

Create a [Secret](https://kubernetes.io/docs/concepts/configuration/secret/)
to hold your git provider access credentials.
- [HTTPS repositories](#https-repositories)
- [SSH repositories](#ssh-repositories)


And add the name to the relevant repo in the `spec`:

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileCatalogSource
metadata:
  name: nginx-catalog
spec:
  repositories:
  - url: https://github.com/weaveworks/profiles-examples
    secretRef:
      name: <name of secret>
```

#### HTTPS repositories

The secret must contain `username` and `password`
fields.

E.g.:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: catalog-secret-https
data:
  username: kewl-beans-1
  password: 1234abcd
```

#### SSH repositories

The secret must contain `identity`, `identity.pub` and
`known_hosts` fields.

E.g.:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: catalog-secret-ssh
data:
  identity: # private ssh key
  identity_pub: # public ssh key
  known_hosts: # known git provider hosts
```

To generate a secret for an SSH repository you can use the following example `kubectl`
command, in this case for GitHub:

```bash
ssh-keygen -q -N "" -f ./identity
ssh-keyscan github.com > ./known_hosts

kubectl create secret generic ssh-credentials \
    --from-file=./identity \
    --from-file=./identity.pub \
    --from-file=./known_hosts

# if your private key is password protected, add
    --from-literal=password=<passphrase>
```

## Create a manual catalog source

Catalog operators also have the option of manually declaring their catalog entries.

The following example will add a single profile to the catalog:

```yaml
apiVersion: weave.works/v1alpha1
kind: ProfileCatalogSource
metadata:
  name: nginx-catalog
spec:
  profiles:
    - name: nginx
      description: This installs some nginx.
      tag: v1.0.0
      url: https://github.com/weaveworks/nginx-profile
      maintainer: weaveworks (https://github.com/weaveworks/profiles)
```

The source manifest can be applied the same way.

```bash
$ kubectl apply -f manual-catalog-source.yaml
# in this case the addition will be instantaneous

$ pctl get --catalog
CATALOG/PROFILE         VERSION DESCRIPTION
nginx-catalog/nginx     v1.0.0  This installs some nginx.
```

## Updating catalog sources

Catalog sources can be updated in the same way as other Kubernetes resources.
Simply edit the manifest and `apply` the changes.

## Removing profiles from the catalog

Likewise, removing a catalog source, and its profiles, is also straightforward:

```sh
$ kubectl delete -f manual-catalog-source.yaml
$ kubectl delete -f dynamic-catalog-source.yaml
$ pctl get --catalog
✗ No available catalog profiles found.
```
