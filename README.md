# profiles
Gitops native package management

## Installation on a local cluster using [Kind](https://kind.sigs.k8s.io/)
1. Set up local environment: `make local-env`.

    This will start a local `kind` cluster and installs
    the `profilesubscription`, `source` and `helm` controllers.

1. deploy an example catalog `kubectl apply -f examples/profile-catalog-source.yaml`

1. To query the catalog API run `kubectl -n profiles-system port-forward <profiles-controller-pod-name> 8000:8000` to enable access to the API and use
[pctl](https://github.com/weaveworks/pctl) to query, for example: `pctl --catalog-url http://localhost:8000 show <search-string>`

1. To see more details on a specific Profile in the catalog, use
[pctl](https://github.com/weaveworks/pctl): `pctl --catalog-url http://localhost:8000 show <profile-name>`

1. Currently `pctl` does not support creating the profile subscription resource from the catalog for you, use the example resource `examples/profile-subscription.yaml` to
subscribe to the example [nginx-profile](https://github.com/weaveworks/nginx-profile): `kubectl apply -f examples/profile-subscription.yaml`

1. The following resources will be created as part of a Helm-based Profile install:
    - ProfileSubscription (the parent object)
    - Profile (the definition as pulled from the upstream target profile)
    - HelmRelease (wrapper resource around the chart)
    - GitRepository (reference to the location of the chart)

    To check your subscription was successful, you can inspect the nginx pod:
    `kubectl describe pod [-n <namespace>] <pod-name>`.
    The pod name will be comprised of `profileSubscriptionName-profileDefinitionName-artifactName-xxxx`

