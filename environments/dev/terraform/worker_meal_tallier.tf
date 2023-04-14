resource "google_project_iam_custom_role" "meal_plan_tallier_role" {
  role_id     = "meal_plan_tallier_role"
  title       = "Meal Plan tallier role"
  description = "An IAM role for the Meal Plan tallier"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
    "pubsub.subscriptions.consume",
    "pubsub.subscriptions.create",
    "pubsub.subscriptions.delete",
  ]
}

locals {
  meal_plan_tallier_database_username = "meal_plan_tallier_db_user"
}

resource "google_pubsub_topic" "meal_plan_tallying_work_topic" {
  name = "meal_plan_tallying_work"
}

resource "google_storage_bucket" "meal_plan_tallier_bucket" {
  name                        = "meal-plan-tallier-cloud-function"
  location                    = "US"
  uniform_bucket_level_access = true
}

data "archive_file" "meal_plan_tallier_function" {
  type        = "zip"
  source_dir  = "${path.module}/meal_plan_tallier_cloud_function"
  output_path = "${path.module}/meal_plan_tallier_cloud_function.zip"
}

resource "google_storage_bucket_object" "meal_plan_tallier_archive" {
  name   = format("meal_plan_tallier_function-%s.zip", data.archive_file.meal_plan_tallier_function.output_md5)
  bucket = google_storage_bucket.meal_plan_tallier_bucket.name
  source = "${path.module}/meal_plan_tallier_cloud_function.zip"
}

resource "google_service_account" "meal_plan_tallier_user_service_account" {
  account_id   = "meal-plan-tallier-worker"
  display_name = "Meal Plans Tallier"
}

resource "random_password" "meal_plan_tallier_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_secret_manager_secret" "meal_plan_tallier_user_database_password" {
  secret_id = "meal_plan_tallier_user_database_password"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "meal_plan_tallier_user_database_password" {
  secret = google_secret_manager_secret.meal_plan_tallier_user_database_password.id

  secret_data = random_password.meal_plan_tallier_user_database_password.result
}

resource "google_sql_user" "meal_plan_tallier_user" {
  name     = local.meal_plan_tallier_database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.meal_plan_tallier_user_database_password.result
}

resource "google_project_iam_member" "meal_plan_tallier_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.meal_plan_tallier_role.id
  member  = format("serviceAccount:%s", google_service_account.meal_plan_tallier_user_service_account.email)
}

# Permissions on the service account used by the function and Eventarc trigger
resource "google_project_iam_member" "meal_plan_tallier_invoking" {
  project = local.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.meal_plan_tallier_user_service_account.email}"
}

resource "google_project_iam_member" "meal_plan_tallier_event_receiving" {
  project    = local.project_id
  role       = "roles/eventarc.eventReceiver"
  member     = "serviceAccount:${google_service_account.meal_plan_tallier_user_service_account.email}"
  depends_on = [google_project_iam_member.meal_plan_tallier_invoking]
}

resource "google_project_iam_member" "meal_plan_tallier_artifactregistry_reader" {
  project    = local.project_id
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${google_service_account.meal_plan_tallier_user_service_account.email}"
  depends_on = [google_project_iam_member.meal_plan_tallier_event_receiving]
}

resource "google_cloudfunctions2_function" "meal_plan_tallier" {
  depends_on = [
    google_project_iam_member.meal_plan_tallier_event_receiving,
    google_project_iam_member.meal_plan_tallier_artifactregistry_reader,
  ]

  name        = "meal-plan-tallier"
  description = "Meal Plan Tallier"
  location    = local.gcp_region

  build_config {
    runtime     = local.go_runtime
    entry_point = "TallyMealPlan"

    source {
      storage_source {
        bucket = google_storage_bucket.meal_plan_tallier_bucket.name
        object = google_storage_bucket_object.meal_plan_tallier_archive.name
      }
    }
  }

  service_config {
    available_memory               = "128Mi"
    ingress_settings               = "ALLOW_INTERNAL_ONLY"
    all_traffic_on_latest_revision = true
    service_account_email          = google_service_account.meal_plan_tallier_user_service_account.email

    environment_variables = {
      # TODO: use the meal_plan_tallier_user for this, currently it has permission denied for accessing tables
      # https://dba.stackexchange.com/questions/53914/permission-denied-for-relation-table
      # https://www.postgresql.org/docs/13/sql-alterdefaultprivileges.html
      PRIXFIXE_DATABASE_USER                     = google_sql_user.api_user.name,
      PRIXFIXE_DATABASE_NAME                     = local.database_name,
      PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME = google_sql_database_instance.dev.connection_name,
      GOOGLE_CLOUD_SECRET_STORE_PREFIX           = format("projects/%d/secrets", data.google_project.project.number)
      GOOGLE_CLOUD_PROJECT_ID                    = data.google_project.project.project_id
    }

    secret_environment_variables {
      key        = "PRIXFIXE_DATABASE_PASSWORD"
      project_id = local.project_id
      secret     = google_secret_manager_secret.api_user_database_password.secret_id
      version    = "latest"
    }
  }

  event_trigger {
    trigger_region        = local.gcp_region
    event_type            = local.pubsub_topic_publish_event
    pubsub_topic          = google_pubsub_topic.meal_plan_tallying_work_topic.id
    retry_policy          = "RETRY_POLICY_RETRY"
    service_account_email = google_service_account.meal_plan_tallier_user_service_account.email
  }
}