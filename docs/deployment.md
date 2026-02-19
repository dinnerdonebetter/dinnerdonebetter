# Deployment Guide

## Branch Strategy

- **Default branch**: `dev` — clone and work from `dev` by default.
- **Deployment branch**: `prod` — production deploys are triggered by pushing to `prod`.

To deploy to production, open a pull request from `dev` into `prod`, then merge it. Merging the PR pushes to `prod` and triggers the deployment pipeline.

**Important**: Never merge directly into `prod` from any branch other than `dev`. All production changes must first land in `dev`.

### Enforcing “prod only from dev”

This repo includes [`.github/workflows/require_dev_for_prod.yaml`](.github/workflows/require_dev_for_prod.yaml), which runs on PRs targeting `prod` and fails unless the source branch is `dev`. In **Settings → Branches → Branch protection rules** for `prod`, add the "check-source-branch" job as a required status check.

---

## What Happens on Deploy

Pushing to `prod` triggers [`.github/workflows/deploy_prod.yaml`](.github/workflows/deploy_prod.yaml). The pipeline runs three jobs in sequence:

| Job                | Description                                                                                           |
|--------------------|-------------------------------------------------------------------------------------------------------|
| **baseline-infra** | Terraform for GCP baseline (project, GKE cluster, networking, DNS, email, IAM)                        |
| **backend-infra**  | Terraform for backend resources (Cloud SQL, storage buckets, Pub/Sub, Algolia, Kubernetes namespaces) |
| **applications**   | Builds and deploys app workloads via Skaffold (API server, workers, cronjobs)                         |

Infrastructure is split into:

- **`infra/deploy/environments/prod/`** — GKE cluster, networking (Cloudflare), email (SendGrid), ingress, cert-manager, external-dns
- **`backend/deploy/environments/prod/`** — Database, storage, Pub/Sub, search, and application Terraform/Kustomize

---

## Local Production Deploy (Applications Only)

To run the application deploy locally (e.g., for debugging or manual rollouts), use:

```bash
./scripts/deploy-prod-local.sh
```

**Prerequisites**: Docker, gcloud, kubectl, skaffold. Your `kubectl` context must point at the prod GKE cluster.

This script only runs the Skaffold application deploy. It does **not** run Terraform or update infrastructure.

---

## Configuration and Secrets

Backend configuration and secrets are documented in [backend/docs/configuration.md](backend/docs/configuration.md). Terraform uses Terraform Cloud for variables; application secrets are synced from GCP Secret Manager into the cluster.
