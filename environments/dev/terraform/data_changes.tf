resource "google_pubsub_topic" "data_changes_topic" {
  name = "data_changes"
}

data "archive_file" "dummy_zip" {
  type        = "zip"
  output_path = "${path.module}/data_changes_function.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "data-changes-function"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "data_changes_function.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.dummy_zip.output_path
}

resource "google_cloudfunctions_function" "data_changes" {
  name        = "data-changes"
  description = "Data Changes"
  runtime     = local.go_runtime

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.data_changes_topic.name
  }

  entry_point = "ProcessDataChange"
}