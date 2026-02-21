#!/usr/bin/env bash
set -euo pipefail

# Proxy Cloud SQL database
# Usage: proxy_db.sh [gcp_project_name] [instance_name]

GCP_PROJECT_NAME="${1:-dinner-done-better-prod}"
INSTANCE_NAME="${2:-prod}"

cloud_sql_proxy "${GCP_PROJECT_NAME}:us-central1:${INSTANCE_NAME}" --port 5434 --gcloud-auth
