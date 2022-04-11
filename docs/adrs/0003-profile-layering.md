# 3. Profile Layering

Date: 2021-11-16

## Status

Accepted

## Context

There are cases where we require to install profiles in a specific order,
perhaps the components being installed are required to be functional before
further workloads are installed into the cluster.

Flux HelmReleases support this through the [`dependsOn` mechanism](https://fluxcd.io/docs/components/helm/helmreleases/#helmrelease-dependencies).

Because the `HelmReleases` are generated inside the CAPI service, we want to
provide a way to configure this dependency without having to manually edit the
files.

We may not always know the specific dependency to be able to set an explicit
dependency, for example, "certificate management" should be in place before
"service mesh".

When installing a set of profiles, it should be possible to declaratively say
that some profiles should be installed before others if this is required.

## Decision

Support Helm Chart [annotation](https://helm.sh/docs/topics/charts/#the-chartyaml-file) that indicates a layer-like ordering for unrelated profiles.

The `weave.works/layer` indicates that a Profile chart should be applied in a
specific layer.

The idea is that we'd order the selected profiles by their layer, and then setup
the HelmRelease dependencies based on the layer ordering.

For example, with the following two Profile charts.


```yaml
apiVersion: v2
name: base-profile
version: 1.0.3
annotations:
  weave.works/profile: base
  weave.works/layer: layer-0
```

```yaml
apiVersion: v2
name: demo-profile
version: 0.1.1
annotations:
  weave.works/profile: demo
  weave.works/layer: layer-1
```

This would result in HelmReleases that looked like this:


```yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: base-profile-default
spec:
  chart:
    spec:
      chart: base-profile
      version: "1.0.3"
      sourceRef:
        kind: HelmRepository
        name: demo
```

With the dependent HelmRelease:

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: demo-profile-default
  namespace: default
spec:
  chart:
    spec:
      chart: demo-profile
      version: "0.1.1"
      sourceRef:
        kind: HelmRepository
        name: demo
  dependsOn:
    - name: base-profile-default
```

The layers are sorted lexicographically, which provides flexibility in naming,
but we can determine a recommended layer naming strategy for co-ordinated use.

The names of the layers are not significant, other than for determining the
ordering of the dependencies.

Given a set of profiles with layers:

 * layer-0
 * layer-1
 * layer-2

This would result in installation with dependencies:

 layer-2 profiles _depend on_ layer-1 profiles, which _depend on_ layer-0
profiles.

## Profiles without layers

Profiles without explicit layers should be configured to depend on the
last-to-be-applied layer.

In the above example, any profiles being installed with no layer would be
configured to depend on layer-2.

## Alternatives

If we don't implement this feature, users can edit the generated `HelmRelease` objects and configure
the correct dependencies at the PR stage.

## Consequences

Installing profiles with dependent profiles into newly bootstrapped clusters is
easier, and more reliable.

The layers mechanism is restricted to the scope of the API call to bootstrap the
cluster, there is nothing that prevents someone changing the resources at the PR
stage tho'.

Lack of co-ordination could lead to confusing layer naming (and possibly broken
installations).

This doesn't preclude a more explicit dependency mechanism.
