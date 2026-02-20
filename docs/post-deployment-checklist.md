# Post-Deployment Checklist

Run this checklist after merging to `prod` and deploying. Items marked with a script idea can be automated.

---

## 1. Public Endpoints (TLS + DNS)

Verify each domain resolves and responds over HTTPS from the public internet.

- [ ] **api.dinnerdonebetter.com** (gRPC, port 443)
  - `curl -sI https://api.dinnerdonebetter.com` or `grpcurl -insecure api.dinnerdonebetter.com:443 list`
- [ ] **http-api.dinnerdonebetter.com** (REST/HTTP)
  - `curl -sI https://http-api.dinnerdonebetter.com/_ops_/ready`
- [ ] **admin.dinnerdonebetter.com** (Admin webapp)
  - `curl -sI https://admin.dinnerdonebetter.com/login`

**Script idea:** `scripts/post-deploy-check-endpoints.sh` — loop over URLs, `curl -sf --connect-timeout 5`, exit 1 if any fail.

---

## 2. TLS Certificate

- [ ] **cert-manager certificate is Ready**
  - `kubectl get certificate -n prod dinner-done-better-cert`
  - Status should be `Ready=True`
- [ ] **Secret has cert data**
  - `kubectl get secret -n prod dinner-done-better-cert -o jsonpath='{.data.tls\.crt}' | base64 -d | head -1`
  - Should start with `-----BEGIN CERTIFICATE-----`
- [ ] **No rate limit errors**
  - `kubectl describe certificate -n prod dinner-done-better-cert`
  - No `429 rateLimited` in Conditions/Events

**Monitor idea:** Grafana alert if `dinner-done-better-cert` Ready=False or renewal is near failure.

---

## 3. Ingress & Load Balancers

- [ ] **Main ingress has an Address**
  - `kubectl get ingress -n prod dinner-done-better-ingress`
  - `ADDRESS` column should show an IP (not blank)
- [ ] **No translation/sync errors**
  - `kubectl describe ingress -n prod dinner-done-better-ingress`
  - No `Translation failed` or `Error 404: sslCertificates`
- [ ] **Admin ingress** (if using standalone)
  - `kubectl get ingress -n prod dinner-done-better-admin-ingress`
  - Should have Address if main ingress is broken

---

## 4. Workload Health (Kubernetes)

- [ ] **API server** — `dinner-done-better-service-api-deployment`
  - `kubectl get pods -n prod -l app.kubernetes.io/name=dinner-done-better-backend`
  - 1/1 Running, no CrashLoopBackOff
- [ ] **Admin webapp** — `dinner-done-better-admin-webapp-deployment`
  - `kubectl get pods -n prod -l app.kubernetes.io/name=dinner-done-better-admin-webapp`
- [ ] **Async message handler** — `dinner-done-better-async-message-handler-deployment`
  - `kubectl get pods -n prod -l app.kubernetes.io/name=dinner-done-better-async-message-handler`
- [ ] **OpenTelemetry collector**
  - `kubectl get pods -n prod -l app.kubernetes.io/name=opentelemetry-collector`

**Script idea:** `kubectl get pods -n prod -o jsonpath='{range .items[?(@.status.phase!="Running")]}{.metadata.name}{"\n"}{end}'` — fail if any non-Running.

---

## 5. CronJobs

- [ ] **All CronJobs exist and have recent last schedule**
  - `kubectl get cronjobs -n prod`
  - Jobs: `dinner-done-better-job-db-cleaner`, `meal-plan-finalizer`, `meal-plan-grocery-list-init`, `meal-plan-task-creator`, `search-data-index-scheduler`
- [ ] **No stuck/failed Job pods**
  - `kubectl get jobs -n prod`
  - `kubectl get pods -n prod -l job-name --field-selector=status.phase!=Succeeded,status.phase!=Running`

**Monitor idea:** Alert if a CronJob has no successful run in the last expected interval (e.g. db-cleaner daily).

---

## 6. Observability (Grafana Cloud)

- [ ] **Metrics in Grafana**
  - Grafana Cloud → Prometheus — verify `dinner_done_better` (or service names) data arriving
  - Note: `prometheusremotewrite` may be commented out in otel config; enable if metrics are missing
- [ ] **Logs in Grafana**
  - Grafana Cloud → Loki — verify logs from API, async handler, admin, cronjobs
  - Look for `service_name` = `dinner_done_better`, `admin_webapp`, etc.
- [ ] **Traces in Grafana**
  - Grafana Cloud → Tempo — verify traces for recent requests

**Monitor idea:** Grafana synthetic check that a trace appears for a health-check request to `http-api.dinnerdonebetter.com/_ops_/ready`.

---

## 7. Pub/Sub Queues

- [ ] **No messages stuck in dead letter topics**
  - GCP Console → Pub/Sub → Subscriptions
  - Check: `data_changes`, `outbound_emails`, `search_index_requests`, `user_data_aggregation_requests`, `webhook_execution_requests`
  - Dead letter topics: `*-deadletter` — ideally 0 unacked
- [ ] **Subscription backlogs reasonable**
  - `gcloud pubsub subscriptions pull --auto-ack <sub-id>` or check metrics in GCP
- [ ] **Async handler consuming**
  - Trigger a flow that publishes (e.g. user signup → outbound_emails), confirm delivery

**Script idea:** `gcloud pubsub subscriptions describe` + parse `numUndeliveredMessages`; alert if > threshold.

---

## 8. External Services

- [ ] **SendGrid** — domain auth, deliverability
  - Send a test email (e.g. password reset or welcome) and confirm inbox delivery
- [ ] **Algolia** — indices populated
  - Algolia Dashboard → Indices: `recipes`, `meals`, `valid_ingredients`, etc.
  - Run a search in the app and confirm results
- [ ] **Cloud SQL (Postgres)** — connectivity from workloads
  - Pods use private IP; no direct curl. Implicit if API responds.

---

## 9. Database & Migrations

- [ ] **Migrations applied**
  - API server runs migrations on startup; successful health check implies they ran
  - Check API logs for migration errors: `kubectl logs -n prod -l app.kubernetes.io/name=dinner-done-better-backend --tail=200`

---

## 10. Secrets & Config

- [ ] **api-service-secrets** exists and is populated
  - `kubectl get secret -n prod api-service-secrets`
- [ ] **admin-webapp-config** (secret + configmap) exist
  - `kubectl get secret -n prod admin-webapp-config`
  - `kubectl get configmap -n prod admin-webapp-config`
- [ ] **grafana-cloud-creds** for otel-collector
  - `kubectl get secret -n prod grafana-cloud-creds`

---

## Automation Opportunities

| Check | Automation |
|-------|------------|
| Endpoint reachability | Bash script + cron or GitHub Action |
| Cert status | `kubectl` + Grafana/k8s alert |
| Pod status | K8s liveness/readiness + Prometheus `kube_pod_status_phase` |
| CronJob success | Grafana alert on `kube_job_status_succeeded` |
| Dead letter queue depth | GCP Monitoring alert on Pub/Sub metrics |
| Grafana data flow | Synthetic trace + Loki log check |

---

## Quick One-Liner (Manual Sanity Check)

```bash
# HTTP endpoints (http-api, admin)
for url in https://http-api.dinnerdonebetter.com/_ops_/ready https://admin.dinnerdonebetter.com/login; do
  echo -n "$url: "; curl -sf -o /dev/null -w "%{http_code}\n" --connect-timeout 5 "$url" || echo "FAIL"
done
# gRPC (api) - requires grpcurl: grpcurl -connect-timeout 5 api.dinnerdonebetter.com:443 list
```
