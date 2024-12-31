resource "google_pubsub_topic" "data_changes_topic" {
  name = "data_changes"
}

resource "google_pubsub_topic" "outbound_emails_topic" {
  name = "outbound_emails"
}

resource "google_pubsub_topic" "search_index_requests_topic" {
  name = "search_index_requests"
}

resource "google_pubsub_topic" "user_data_aggregator_topic" {
  name = "user_data_aggregation_requests"
}

resource "google_pubsub_topic" "webhook_execution_requests_topic" {
  name = "webhook_execution_requests"
}
