resource "google_pubsub_topic" "data_changes_topic" {
  name = "data_changes"
}

resource "google_storage_bucket" "data_changes_cloud_function_bucket" {
  name     = "data-changes-cloud-function"
  location = "US"
}

resource "google_cloudfunctions_function" "data_changes" {
  name        = "data-changes"
  description = "Data Changes"
  runtime     = local.go_runtime

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.data_changes_cloud_function_bucket.name

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.data_changes_topic.name
  }

  entry_point = "ProcessDataChange"
}