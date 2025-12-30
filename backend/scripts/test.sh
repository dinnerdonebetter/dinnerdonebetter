#!/usr/bin/env bash
set -euo pipefail

# Run tests
# Usage: test.sh <go_test_command> <testable_package_list>

TESTABLE_PACKAGE_LIST="${1}"

# shellcheck disable=SC2086
CGO_ENABLED=1 go test -shuffle=on -race -vet=all -failfast ${TESTABLE_PACKAGE_LIST}

