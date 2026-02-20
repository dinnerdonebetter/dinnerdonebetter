# Production Cost-Saving Opportunities

Analysis of the production deployment (Terraform/Kustomize/Skaffold) for low-effort cost optimizations.

---

## 1. Consolidate the Two GCE Ingresses (Biggest Win)

You have two separate Ingresses:

- `infra/deploy/environments/prod/kustomize/ingress.yaml` — api, http-api, admin, app, www (TLS hosts)
- `infra/deploy/environments/prod/kustomize/ingress-admin.yaml` — admin only

Each GCE Ingress creates its own load balancer (~$18/month per forwarding rule). Merging into one removes one load balancer.

**Est. savings:** ~$18/month

**Steps:**
1. Add the admin rule from `ingress-admin.yaml` into `ingress.yaml` (admin may already be in the main ingress TLS; ensure there's a rule for it).
2. Remove `ingress-admin.yaml` from `kustomization.yaml` resources.
3. Delete the `ingress-admin.yaml` file.
4. Deploy and verify admin.dinnerdonebetter.com resolves and works.
5. (The admin ingress was split due to "transient GCE controller issues"; if merging causes problems, revert.)

---

## 2. Do Not Add More Ingresses

Each new GCE Ingress creates a new load balancer (~$18/month). When troubleshooting (e.g., gRPC backend stuck UNHEALTHY), avoid splitting hosts into separate Ingresses. Instead:

- Ensure `BackendConfig` (health check via HTTP port) and `cloud.google.com/app-protocols: '{"grpc":"HTTP2"}'` are on the API service for the gRPC port.
- If a backend remains stuck, try: remove the api rule from the main Ingress, `kubectl apply`, wait 5–10 min for GCE to garbage-collect the backend, add the rule back, apply again to create a fresh backend.

---

## 3. Orphaned Static IP

`infra/deploy/environments/prod/terraform/networking.tf` defines `google_compute_address.static_ip` ("prod") but it's not referenced anywhere (no Ingress annotation, etc.). If reserved and unattached, it costs ~$7/month.

**Est. savings:** ~$7/month

**Steps:**
1. Verify the static IP isn't used by anything (check GCP Console → VPC network → IP addresses).
2. If unused, remove the `google_compute_address` resource from `networking.tf`.
3. Run `terraform plan` and `terraform apply`.

---

## 4. Optional: Cloud SQL Public IP

Cloud SQL has `ipv4_enabled = true` and a Cloudflare A record for `db.dinnerdonebetter.com`. If all access is via private VPC (e.g., GKE), consider disabling the public IP for security and minor cost/simplification. Only do this if nothing needs direct public DB access.

**Difficulty:** Medium — verify no external tools/CI rely on public DB access first.

---

## Summary

| Opportunity                      | Est. monthly savings | Difficulty |
|----------------------------------|----------------------|------------|
| Merge 2 Ingresses into 1         | ~$18                 | Low        |
| Don't add more Ingresses (gRPC)  | ~$18 avoided          | —          |
| Remove orphaned static IP        | ~$7                  | Low        |
| **Total low-hanging fruit**      | **~$25/month**       |            |
