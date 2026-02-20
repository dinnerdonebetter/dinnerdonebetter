#!/usr/bin/env bash
# Post-deploy verification for prod.
# Run after skaffold deploy. Requires: kubectl (configured for prod cluster), curl, grpcurl.
set -euo pipefail

NAMESPACE="${NAMESPACE:-prod}"
CONNECT_TIMEOUT="${CONNECT_TIMEOUT:-15}"
MAX_RETRIES="${MAX_RETRIES:-3}"
RETRY_DELAY="${RETRY_DELAY:-10}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

fail() {
  echo "FAIL: $*"
  exit 1
}

ok() {
  echo "OK: $*"
}

retry_cmd() {
  local cmd=("$@")
  local i=1
  while [[ $i -le $MAX_RETRIES ]]; do
    echo "  Attempt $i/$MAX_RETRIES..."
    if "${cmd[@]}"; then
      return 0
    fi
    i=$((i + 1))
    [[ $i -le $MAX_RETRIES ]] && sleep "$RETRY_DELAY"
  done
  return 1
}

echo "=== Post-deploy verification (namespace=$NAMESPACE) ==="

# --- K8s resources ---
echo ""
echo "[1/6] Certificate Ready"
status=$(kubectl get certificate -n "$NAMESPACE" dinner-done-better-cert -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}' 2>/dev/null || true)
if [[ "$status" != "True" ]]; then
  fail "Certificate dinner-done-better-cert not Ready (status=$status)"
fi
ok "Certificate Ready=True"

echo ""
echo "[2/6] Ingress has address"
addr=$(kubectl get ingress -n "$NAMESPACE" dinner-done-better-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || kubectl get ingress -n "$NAMESPACE" dinner-done-better-ingress -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || true)
if [[ -z "${addr:-}" ]]; then
  fail "Ingress dinner-done-better-ingress has no address"
fi
ok "Ingress address: $addr"

echo ""
echo "[3/6] Core workloads Running"
for label in "app.kubernetes.io/name=dinner-done-better-backend" "app.kubernetes.io/name=dinner-done-better-admin-webapp" "app.kubernetes.io/name=dinner-done-better-async-message-handler"; do
  not_running=$(kubectl get pods -n "$NAMESPACE" -l "$label" -o jsonpath='{range .items[?(@.status.phase!="Running")]}{.metadata.name}{"\n"}{end}' 2>/dev/null || true)
  if [[ -n "$not_running" ]]; then
    fail "Pods not Running (label=$label): $not_running"
  fi
done
ok "API, admin, async handler pods Running"

# --- Public endpoints (with retries for LB/DNS propagation) ---
echo ""
echo "[4/6] HTTP API ready"
retry_cmd curl -sf --connect-timeout "$CONNECT_TIMEOUT" --max-time 30 https://http-api.dinnerdonebetter.com/_ops_/ready -o /dev/null || fail "http-api.dinnerdonebetter.com not ready"
ok "http-api.dinnerdonebetter.com/_ops_/ready"

echo ""
echo "[5/6] Admin webapp responds"
check_admin() {
  local code
  code=$(curl -sL -o /dev/null -w "%{http_code}" --connect-timeout "$CONNECT_TIMEOUT" --max-time 30 https://admin.dinnerdonebetter.com/login)
  [[ "${code:0:1}" = "2" ]] || [[ "${code:0:1}" = "3" ]]
}
retry_cmd check_admin || fail "admin.dinnerdonebetter.com not responding (2xx/3xx)"
ok "admin.dinnerdonebetter.com"

echo ""
echo "[6/6] gRPC API list"
if command -v grpcurl >/dev/null 2>&1; then
  retry_cmd grpcurl -insecure -connect-timeout "$CONNECT_TIMEOUT" api.dinnerdonebetter.com:443 list >/dev/null || fail "api.dinnerdonebetter.com gRPC not reachable"
  ok "api.dinnerdonebetter.com gRPC list"
else
  echo "  (skipped: grpcurl not installed; run 'go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest')"
fi

echo ""
echo "=== All verifications passed ==="
