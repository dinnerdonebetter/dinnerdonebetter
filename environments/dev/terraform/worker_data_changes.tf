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

resource "google_cloudfunctions2_function" "data_changes" {
  name        = "data-changes"
  location    = "us-central1"
  description = format("Data Changes (%s)", data.archive_file.data_changes_function.output_md5)

  event_trigger {
    event_type            = local.pubsub_topic_publish_event
    pubsub_topic          = google_pubsub_topic.data_changes_topic.id
    retry_policy          = "RETRY_POLICY_RETRY"
    service_account_email = google_service_account.data_changes_user_service_account.email
  }

  build_config {
    runtime     = local.go_runtime
    entry_point = "ProcessDataChange"

    source {
      storage_source {
        bucket = google_storage_bucket.data_changes_bucket.name
        object = google_storage_bucket_object.data_changes_archive.name
      }
    }
  }

  service_config {
    available_memory = "128M"

    environment_variables = {
      PRIXFIXE_SENDGRID_API_TOKEN      = var.SENDGRID_API_TOKEN
      PRIXFIXE_SEGMENT_API_TOKEN       = var.SEGMENT_API_TOKEN
      GOOGLE_CLOUD_SECRET_STORE_PREFIX = format("projects/%d/secrets", data.google_project.project.number)
      GOOGLE_CLOUD_PROJECT_ID          = data.google_project.project.project_id
    }
  }
}