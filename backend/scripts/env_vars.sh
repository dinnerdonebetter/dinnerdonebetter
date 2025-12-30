#!/usr/bin/env bash
set -euo pipefail

# Generate and format env vars
# Usage: env_vars.sh <package_prefix> <gofmt_command>

PACKAGE_PREFIX="${1:-github.com/dinnerdonebetter/backend}"
GO_FORMAT="${2:-gofmt -s -w}"

go run "${PACKAGE_PREFIX}/cmd/tools/codegen/valid_env_vars"
# shellcheck disable=SC2086
${GO_FORMAT} internal/config/envvars/*.go

