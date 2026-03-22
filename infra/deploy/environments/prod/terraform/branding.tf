locals {
  company_name     = "DinnerDoneBetter"
  company_slug     = "dinner-done-better"
  company_slug_ns  = "dinnerdonebetter"
  public_domain    = "${local.company_slug_ns}.com"
  gcp_project_id   = "${local.company_slug}-prod"
  tf_cloud_org     = local.company_slug_ns
  grafana_prom_ds  = "grafanacloud-${local.company_slug_ns}-prom"
  grafana_loki_ds  = "grafanacloud-${local.company_slug_ns}-logs"
}
