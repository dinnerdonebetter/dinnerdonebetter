resource "google_pubsub_topic" "data_changes_topic" {
  name = "data_changes"
}

resource "google_storage_bucket" "data_changes_bucket" {
  name     = "data-changes-function"
  location = "US"
}

resource "google_storage_bucket_object" "data_changes_archive" {
  name   = "data_changes_function.zip"
  bucket = google_storage_bucket.data_changes_bucket.name
  source = "${path.module}/data_changes_cloud_function.zip"
}

resource "google_service_account" "data_changes_user_service_account" {
  account_id   = "data-changes-worker"
  display_name = "Data Changes Worker"
}

resource "google_project_iam_member" "data_changes_user" {
  project = local.project_id
  role    = "roles/viewer"
  member  = format("serviceAccount:%s", google_service_account.data_changes_user_service_account.email)
}

resource "google_project_iam_binding" "data_changes_user_secret_accessor" {
  project = local.project_id
  role    = "roles/secretmanager.secretAccessor"

  members = [
    google_project_iam_member.data_changes_user.member,
  ]
}

resource "google_project_iam_binding" "data_changes_user_cloud_sql_client" {
  project = local.project_id
  role    = "roles/cloudsql.client"

  members = [
    google_project_iam_member.data_changes_user.member,
  ]
}

resource "google_cloudfunctions_function" "data_changes" {
  name        = "data-changes"
  description = "Data Changes"
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