resource "google_monitoring_uptime_check_config" "api_uptime" {
  display_name = "api-server-uptime-check"
  timeout      = "60s"
  period       = "300s"

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

resource "google_monitoring_notification_channel" "api_server_monitor_notification_channel" {
  display_name = "API Service Monitor Notification Channel"
  type         = "email"

  labels = {
    email_address = "verygoodsoftwarenotvirus@protonmail.com"
  }
}

resource "google_monitoring_alert_policy" "alert_policy" {
  display_name = "Dev API Alert Policy"
  combiner     = "OR"

  notification_channels = [
    google_monitoring_notification_channel.api_server_monitor_notification_channel.name,
  ]

  conditions {
    display_name = "request latency"
    condition_monitoring_query_language {
      duration = ""
      query    = <<END
        fetch uptime_url
        | metric 'monitoring.googleapis.com/uptime_check/request_latency'
        | filter (metric.check_id == 'api-server-uptime-check')
        | group_by 5m, [value_request_latency_max: max(value.request_latency)]
        | every 5m
        | group_by [], [value_request_latency_max_max: max(value_request_latency_max)]
        | group_by [],
            [value_request_latency_max_max_mean: mean(value_request_latency_max_max)]
        | condition val() > 999 'ms'
      END
    }
  }
}
