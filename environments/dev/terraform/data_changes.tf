resource "google_pubsub_topic" "data_changes_topic" {
  name = "data_changes"
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
  role    = "roles/secretmanager.secretAccessor"
  member  = format("serviceAccount:%s", google_service_account.data_changes_user_service_account.email)
}

resource "google_project_iam_binding" "data_changes_user_pubsub_publisher" {
  project = local.project_id
  role    = "roles/pubsub.publisher"

  members = [
    google_project_iam_member.data_changes_user.member,
  ]
}

resource "google_project_iam_binding" "data_changes_user_pubsub_subscriber" {
  project = local.project_id
  role    = "roles/pubsub.subscriber"

  members = [
    google_project_iam_member.data_changes_user.member,
  ]
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