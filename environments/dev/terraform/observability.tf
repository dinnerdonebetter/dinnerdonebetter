resource "google_monitoring_uptime_check_config" "api_uptime" {
  display_name = "api-server-uptime-check"
  timeout      = "60s"

  http_check {
    path         = "/_meta_/ready"
    port         = "443"
    use_ssl      = true
    validate_ssl = true
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = local.project_id
      host       = "api.prixfixe.dev"
    }
  }
}