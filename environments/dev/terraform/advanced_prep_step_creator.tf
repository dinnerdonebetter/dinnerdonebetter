resource "google_project_iam_custom_role" "advanced_prep_step_creator_role" {
  role_id     = "advanced_prep_step_creator_role"
  title       = "Advanced Prep Step Creator Role"
  description = "An IAM role for the advanced prep step creator"
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
  advanced_prep_step_creator_database_username = "advanced_prep_step_creator_db_user"
}

resource "google_pubsub_topic" "advanced_prep_step_creator_topic" {
  name = "advanced_prep_step_creation_work"
}

resource "google_cloud_scheduler_job" "advanced_prep_step_creation" {
  project = local.project_id
  region  = local.gcp_region
  name    = "advanced-prep-step-creation-scheduler"

  schedule  = "* * * * *" # every minute
  time_zone = "America/Chicago"

  pubsub_target {
    topic_name = google_pubsub_topic.advanced_prep_step_creator_topic.id
    data       = base64encode("{}")
  }
}

resource "google_storage_bucket" "advanced_prep_step_creator_bucket" {
  name     = "advanced-prep-step-creation-cloud-function"
  location = "US"
}

data "archive_file" "advanced_prep_step_creator_function" {
  type        = "zip"
  source_dir  = "${path.module}/advanced_prep_step_creator_cloud_function"
  output_path = "${path.module}/advanced_prep_step_creator_cloud_function.zip"
}

resource "google_storage_bucket_object" "advanced_prep_step_creator_archive" {
  name   = format("advanced_prep_step_creator_function-%s.zip", data.archive_file.advanced_prep_step_creator_function.output_md5)
  bucket = google_storage_bucket.advanced_prep_step_creator_bucket.name
  source = "${path.module}/advanced_prep_step_creator_cloud_function.zip"
}

resource "google_service_account" "advanced_prep_step_creator_user_service_account" {
  account_id   = "advanced-prep-step-creation-worker"
  display_name = "Advanced Prep Step Creator"
}

resource "google_project_iam_member" "advanced_prep_step_creator_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.advanced_prep_step_creator_role.id
  member  = format("serviceAccount:%s", google_service_account.advanced_prep_step_creator_user_service_account.email)
}

resource "random_password" "advanced_prep_step_creator_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_secret_manager_secret" "advanced_prep_step_creator_user_database_password" {
  secret_id = "advanced_prep_step_creator_user_database_password"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "advanced_prep_step_creator_user_database_password" {
  secret = google_secret_manager_secret.advanced_prep_step_creator_user_database_password.id

  secret_data = random_password.advanced_prep_step_creator_user_database_password.result
}

resource "google_sql_user" "advanced_prep_step_creator_user" {
  name     = local.advanced_prep_step_creator_database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.advanced_prep_step_creator_user_database_password.result
}

resource "google_cloudfunctions_function" "advanced_prep_step_creator" {
  name                = "advanced-prep-step-creation"
  description         = "Advanced Prep Step Creator"
  runtime             = local.go_runtime
  available_memory_mb = 128

  source_archive_bucket = google_storage_bucket.advanced_prep_step_creator_bucket.name
  source_archive_object = google_storage_bucket_object.advanced_prep_step_creator_archive.name
  service_account_email = google_service_account.advanced_prep_step_creator_user_service_account.email

  entry_point = "CreateAdvancedPrepSteps"

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.advanced_prep_step_creator_topic.name
  }

  environment_variables = {
    # TODO: use the advanced_prep_step_creator_user for this, currently it has permission denied for accessing tables
    # https://dba.stackexchange.com/questions/53914/permission-denied-for-relation-table
    # https://www.postgresql.org/docs/13/sql-alterdefaultprivileges.html
    PRIXFIXE_DATABASE_USER                     = google_sql_user.api_user.name,
    PRIXFIXE_DATABASE_PASSWORD                 = random_password.api_user_database_password.result,
    PRIXFIXE_DATABASE_NAME                     = local.database_name,
    PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME = google_sql_database_instance.dev.connection_name,
    GOOGLE_CLOUD_SECRET_STORE_PREFIX           = format("projects/%d/secrets", data.google_project.project.number)
    GOOGLE_CLOUD_PROJECT_ID                    = data.google_project.project.project_id
  }
}