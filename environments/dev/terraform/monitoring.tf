resource "google_monitoring_alert_policy" "worker_memory_usage_alert_policy" {
  display_name = "Unacked Pub/Sub Messages"
  combiner     = "OR"
  conditions {
    display_name = "Unacked Pub/Sub Messages"

    condition_threshold {
      filter     = "resource.type = \"pubsub_topic\" AND metric.type = \"pubsub.googleapis.com/topic/num_unacked_messages_by_region\""
      duration   = "300s"
      comparison = "COMPARISON_GT"
      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_COUNT"
      }
      trigger {
        count = 1
      }
      threshold_value = 25
    }
  }

  alert_strategy {
    auto_close = "259200s"
  }
}
