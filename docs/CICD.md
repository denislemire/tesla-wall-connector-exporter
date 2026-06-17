# CI/CD (deferred)

**Status:** Not active. Build and publish images locally (`make docker-push`) until this pipeline is wired.

This document captures the intended automation for a self-hosted CircleCI Server environment with in-cluster Kaniko (no VM/Docker executor) and GitOps deploy via Flux.

## Goals

1. Run `go test ./...` on every push and tag.
2. Build container images with **Kaniko** inside the cluster (no machine executor).
3. Push to **ghcr.io** (public) and a **private registry** mirror used by the home cluster.
4. Package and publish the Helm chart (OCI).
5. Let Flux reconcile a `HelmRelease` in the `monitoring` namespace.

## Intended flow

```
Git push / tag
  → CircleCI (OIDC token)
  → Authentik token exchange → kubectl
  → apply Kaniko Job (git clone tag → build → push)
  → helm package + OCI push
  → (optional) Flux ImageUpdateAutomation commits image tag to infra Git
  → Flux HelmRelease upgrade
```

## CircleCI OIDC → Kubernetes

Reference pattern (operator docs in consuming infra repos):

- Job image: `alpine/k8s`
- Context holds: API server URL, cluster CA (base64), Authentik token URL
- Job receives `CIRCLE_OIDC_TOKEN_V2` automatically
- Script exchanges CircleCI JWT for Authentik access token with `aud=kubernetes-cluster`
- RBAC binds group `oidc:circleci-oidc` to permissions needed to create Kaniko Jobs

See `ci/kubeconfig-from-oidc.sh.example` for the token exchange skeleton.

## Kaniko Job (no VM)

Reference: `ci/kaniko-build-job.yaml.example`

- Init container clones the tagged Git release
- Kaniko builds `Dockerfile` and pushes to configured registries
- Registry credentials via Kubernetes Secrets mounted at `/kaniko/.docker`

## Enabling CI later

1. Copy `ci/circleci-continue-config.yml.example` → `.circleci/continue-config.yml` and add a root `config.yml` with path-filtering if desired.
2. Create CircleCI project; add context with cluster/OIDC variables.
3. Ensure Kaniko namespace, registry push secret, and RBAC exist in the cluster.
4. Remove the **deferred** banner from this file once a green pipeline exists.

## Public vs private artifacts

| Artifact | Public | Private mirror |
|----------|--------|----------------|
| Container image | `ghcr.io/denislemire/tesla-wall-connector-exporter:<tag>` | Operator-specific registry path |
| Helm chart | Git tag in this repo | OCI push to operator registry (optional) |
| Flux manifests | N/A | Operator infra Git only |

Nothing in this repo should hard-code operator hostnames or registry URLs beyond the public ghcr.io default in the Helm chart.
