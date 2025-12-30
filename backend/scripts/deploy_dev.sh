#!/usr/bin/env bash
set -euo pipefail

# Deploy to dev environment
# Usage: deploy_dev.sh [skaffold_file] [profile] [build_concurrency]

SKAFFOLD_FILE="${1:-skaffold.yaml}"
PROFILE="${2:-dev}"
BUILD_CONCURRENCY="${3:-3}"

skaffold run --filename="${SKAFFOLD_FILE}" --build-concurrency "${BUILD_CONCURRENCY}" --profile "${PROFILE}"

