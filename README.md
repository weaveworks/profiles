# profiles
Gitops native package management

## Installation on a local cluster using [Kind](https://kind.sigs.k8s.io/)
1. Set up local environment: `make local-env`.

  This will start a local `kind` cluster and installs
  the `profilesubscription`, `source` and `helm` controllers.

1. Then subscribe to the example [nginx-profile](https://github.com/weaveworks/nginx-profile): `kubectl apply -f examples/profile-subscription.yaml`

1. The following resources will be created as part of a Helm-based Profile install:
  - ProfileSubscription (the parent object)
  - Profile (the definition as pulled from the upstream target profile)
  - HelmRelease (wrapper resource around the chart)
  - GitRepository (reference to the location of the chart)

  To check your subscription was successful, you can inspect the nginx pod:
  `kubectl describe pod [-n <namespace>] <pod-name>`.
  The pod name will be comprised of `profileSubscriptionName-profileDefinitionName-artifactName-xxxx`
