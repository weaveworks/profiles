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

## Release process
There are some manual steps right now, should be streamlined soon.

Steps:

1. Create a new release notes file:
	```sh
	touch docs/release_notes/<version>.md
	```

1. Copy-and paste the release notes from the draft on the releases page into this file.
    _Note: sometimes the release drafter is a bit of a pain, verify that the notes are
    correct by doing something like: `git log --first-parent tag1..tag2`._

1. PR the release notes into main.

1. Create and push a tag with the new version:
	```sh
	git tag <version>
	git push origin <version>
	```

1. The `Create release` action should run. Verify that:
	1. The release has been created in Github
		1. With the correct assets
		1. With the correct release notes
	1. The image has been pushed to docker
	1. The image can be pulled and used in a deployment
