resource "google_pubsub_topic" "meal_plan_topic" {
  name = "meal_plan_finalization_work"
}

resource "google_cloud_scheduler_job" "meal_plan_finalization" {
  project = local.project_id
  region  = local.gcp_region
  name    = "meal-plan-finalizer-scheduler"

  schedule  = "* * * * *" # every minute
  time_zone = "America/Chicago"

  pubsub_target {
    topic_name = google_pubsub_topic.meal_plan_topic.id
    data       = base64encode("{}")
  }
}

resource "google_storage_bucket" "meal_plan_finalizer_bucket" {
  name     = "meal-plan-finalizer-cloud-function"
  location = "US"
}

resource "google_storage_bucket_object" "meal_plan_finalizer_archive" {
  name   = "meal_plan_finalizer_function.zip"
  bucket = google_storage_bucket.meal_plan_finalizer_bucket.name
  source = "${path.module}/meal_plan_finalizer_cloud_function.zip"
}

resource "google_service_account" "meal_plan_finalizer_user_service_account" {
  account_id   = "meal-plan-finalizer-worker"
  display_name = "Meal Plans Finalizer"
}

resource "google_project_iam_member" "meal_plan_finalizer_user" {
  project = local.project_id
  role    = "roles/viewer"
  member  = format("serviceAccount:%s", google_service_account.meal_plan_finalizer_user_service_account.email)
}

resource "google_project_iam_binding" "meal_plan_finalizer_user_secret_accessor" {
  project = local.project_id
  role    = "roles/secretmanager.secretAccessor"

  members = [
    google_project_iam_member.meal_plan_finalizer_user.member,
  ]
}

resource "google_project_iam_binding" "meal_plan_finalizer_user_cloud_sql_client" {
  project = local.project_id
  role    = "roles/cloudsql.client"

  members = [
    google_project_iam_member.meal_plan_finalizer_user.member,
  ]
}

resource "google_cloudfunctions_function" "meal_plan_finalizer" {
  name        = "meal-plan-finalizer"
  description = "Meal Plan Finalizer"
  runtime     = local.go_runtime

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.meal_plan_finalizer_bucket.name
  source_archive_object = google_storage_bucket_object.meal_plan_finalizer_archive.name
  service_account_email = google_service_account.meal_plan_finalizer_user_service_account.email

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.meal_plan_topic.name
  }

  entry_point = "FinalizeMealPlans"
}