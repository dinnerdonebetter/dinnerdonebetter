resource "google_pubsub_topic" "data_changes_topic" {
  name = "data_changes"
}

resource "google_project_iam_custom_role" "data_changes_worker_role" {
  role_id     = "data_changes_worker_role"
  title       = "Data changes worker role"
  description = "An IAM role for the data changes worker"
  permissions = [
    "secretmanager.versions.access",
    "pubsub.topics.list",
    "pubsub.subscriptions.consume",
    "pubsub.subscriptions.create",
    "pubsub.subscriptions.delete",
  ]
}

resource "google_storage_bucket" "data_changes_bucket" {
  name     = "data-changes-cloud-function"
  location = "US"
}

data "archive_file" "data_changes_function" {
  type        = "zip"
  source_dir  = "${path.module}/data_changes_cloud_function"
  output_path = "${path.module}/data_changes_cloud_function.zip"
}

resource "google_storage_bucket_object" "data_changes_archive" {
  name   = format("data_changes_function-%s.zip", data.archive_file.data_changes_function.output_md5)
  bucket = google_storage_bucket.data_changes_bucket.name
  source = "${path.module}/data_changes_cloud_function.zip"
}

resource "google_service_account" "data_changes_user_service_account" {
  account_id   = "data-changes-worker"
  display_name = "Data Changes Worker"
}

resource "google_project_iam_member" "data_changes_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.data_changes_worker_role.id
  member  = format("serviceAccount:%s", google_service_account.data_changes_user_service_account.email)
}

resource "google_cloudfunctions_function" "data_changes" {
  name        = "data-changes"
  description = format("Data Changes (%s)", data.archive_file.data_changes_function.output_md5)
  runtime     = local.go_runtime

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.data_changes_bucket.name
  source_archive_object = google_storage_bucket_object.data_changes_archive.name
  service_account_email = google_service_account.data_changes_user_service_account.email

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.data_changes_topic.name
  }

  entry_point = "ProcessDataChange"
}