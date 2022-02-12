resource "google_project_iam_custom_role" "meal_plan_prober_role" {
  role_id     = "meal_plan_prober_role"
  title       = "Meal Plan prober role"
  description = "An IAM role for the meal plan prober"
  permissions = [
    "pubsub.subscriptions.consume",
  ]
}

resource "google_pubsub_topic" "meal_plan_prober_topic" {
  name = "meal_plan_prober_work"
}

resource "google_cloud_scheduler_job" "meal_plan_prober_schedule" {
  project = local.project_id
  region  = local.gcp_region
  name    = "meal-plan-prober-scheduler"

  schedule  = "*/5 * * * *" # every five minutes
  time_zone = "America/Chicago"

  pubsub_target {
    topic_name = google_pubsub_topic.meal_plan_prober_topic.id
    data       = base64encode("{}")
  }
}

resource "google_storage_bucket" "meal_plan_prober_bucket" {
  name     = "meal-plan-prober-cloud-function"
  location = "US"
}

data "archive_file" "meal_plan_prober_function" {
  type        = "zip"
  source_dir  = "${path.module}/meal_planning_prober_cloud_function"
  output_path = "${path.module}/meal_planning_prober_cloud_function.zip"
}

resource "google_storage_bucket_object" "meal_plan_prober_archive" {
  name   = format("meal_planning_prober_function-%s.zip", data.archive_file.meal_plan_prober_function.output_md5)
  bucket = google_storage_bucket.meal_plan_prober_bucket.name
  source = "${path.module}/meal_planning_prober_cloud_function.zip"
}

resource "google_service_account" "meal_plan_prober_user_service_account" {
  account_id   = "meal-planning-prober-worker"
  display_name = "Meal Planning Prober"
}

resource "google_project_iam_member" "meal_plan_prober_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.meal_plan_prober_role.id
  member  = format("serviceAccount:%s", google_service_account.meal_plan_prober_user_service_account.email)
}

resource "google_cloudfunctions_function" "meal_plan_prober" {
  name                = "meal-planning-prober"
  description         = "Meal Planning Prober"
  runtime             = local.go_runtime
  available_memory_mb = 128

  source_archive_bucket = google_storage_bucket.meal_plan_prober_bucket.name
  source_archive_object = google_storage_bucket_object.meal_plan_prober_archive.name
  service_account_email = google_service_account.meal_plan_prober_user_service_account.email

  entry_point = "ProbeMealPlanning"

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.meal_plan_prober_topic.name
  }

  environment_variables = {
    # TODO: use the meal_plan_prober_user for this, currently it has permission denied for accessing tables
    # https://dba.stackexchange.com/questions/53914/permission-denied-for-relation-table
    # https://www.postgresql.org/docs/13/sql-alterdefaultprivileges.html
    PRIXFIXE_DATABASE_USER                     = google_sql_user.api_user.name,
    PRIXFIXE_DATABASE_PASSWORD                 = random_password.api_user_database_password.result,
    PRIXFIXE_DATABASE_NAME                     = local.database_name,
    PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME = google_sql_database_instance.dev.connection_name,
    GOOGLE_CLOUD_SECRET_STORE_PREFIX           = format("projects/%d/secrets", data.google_project.project.number)
  }

  #  secret_environment_variables = {
  #    key    = "PRIXFIXE_DATA_CHANGES_TOPIC_NAME"
  #    secret = google_secret_manager_secret.data_changes_topic_name.id
  #  }
  #
  #  secret_volumes = {
  #    mount_path = "/config/"
  #    secret     = google_secret_manager_secret.api_service_config.id
  #    versions = {
  #      path = "config.json"
  #    }
  #  }
}