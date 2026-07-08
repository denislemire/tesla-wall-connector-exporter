# CI/CD (EhWS CircleCI)

**Status:** Active on CircleCI Server (`circle.ehws.generic.business`) via KubeVirt machine executor.

## Goals

1. Build container images with **Docker on a KubeVirt `linux.medium` VM** (ubuntu-2404 golden image).
2. Push to the **EhWS Zot registry** at `registry.ehws.generic.business`.
3. Let Flux reconcile the `HelmRelease` in `ehws-infra` (`clusters/ehws/twc-exporter`).

Public **ghcr.io** images remain a separate local or GitHub Actions path (`make docker-push`).

## Flow

```
Git push (main / feat branch) or v* tag
  → CircleCI machine job (KubeVirt ubuntu-2404)
  → 1Password Secrets Automation (ehws-1password-secrets context)
  → docker login registry.ehws.generic.business
  → docker build + push (multi-tag on release)
  → Flux HelmRelease picks up new tag in ehws-infra
```

## CircleCI setup

| Item | Value |
|------|-------|
| Config | `.circleci/config.yml` |
| Executor | `machine` / `ubuntu-2404:current` / `resource_class: medium` |
| Context | `ehws-1password-secrets` (`OP_SERVICE_ACCOUNT_TOKEN`) |
| Orb | `onepassword/secrets@1` |
| Zot creds | `op://EhWS Secrets/Registry EHWS - Zot htpasswd/{username,password}` |

See `ehws-infra` → `docs/CIRCLECI_1PASSWORD_SECRETS.md` for context and 1Password service-account setup.

### Image tags

| Trigger | Tags pushed |
|---------|-------------|
| Branch push (`main`, `feat/circleci-kubevirt-zot-build`) | `sha-<7>`, `branch-<name>` |
| Git tag `v*` | `vX.Y.Z`, `X.Y.Z`, `X.Y`, `sha-<7>` |

Registry path: `registry.ehws.generic.business/tesla-wall-connector-exporter`.

### Deploy (Flux)

After a release tag is pushed, bump `image.tag` in `ehws-infra/clusters/ehws/twc-exporter/helm-release.yaml` and reconcile. Ensure `image.repository` matches the Zot path above.

## KubeVirt / DNS prerequisites

Machine jobs resolve EhWS hostnames via in-cluster CoreDNS static hosts (not LAN split-horizon). Required entries are documented in `ehws-infra/clusters/ehws/circleci-machines/README.md`.

Troubleshooting hangs: `ehws-infra/scripts/diagnose-kubevirt-machine-executor.sh`.

## Legacy reference (Kaniko / OIDC)

The `ci/*.example` files describe an alternate in-cluster Kaniko + OIDC pattern that is **not** the active EhWS pipeline. Kept for reference only.
