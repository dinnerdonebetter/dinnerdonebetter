resource "google_pubsub_topic" "user_data_aggregator_topic" {
  name = "user_data_aggregation_requests"
}

resource "google_project_iam_custom_role" "user_data_aggregator_role" {
  role_id     = "user_data_aggregator_role"
  title       = "user data aggregator role"
  description = "An IAM role for the user data aggregator"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
    "pubsub.subscriptions.consume",
    "pubsub.subscriptions.create",
    "pubsub.subscriptions.delete",
    "eventarc.events.receiveAuditLogWritten",
    "eventarc.events.receiveEvent",
    "run.jobs.run",
    "run.routes.invoke",
    "artifactregistry.dockerimages.get",
    "artifactregistry.dockerimages.list",
    "storage.objects.list",
    "storage.objects.get",
    "storage.objects.update",
    "storage.objects.create",
    "storage.objects.delete",
    "storage.objects.get",
  ]
}

resource "google_storage_bucket" "user_data_aggregator_bucket" {
  name     = "user-data-aggregator-cloud-function"
  location = "US"
}

data "archive_file" "user_data_aggregator_function" {
  type        = "zip"
  source_dir  = "${path.module}/user_data_aggregator_cloud_function"
  output_path = "${path.module}/user_data_aggregator_cloud_function.zip"
}

resource "google_storage_bucket_object" "user_data_aggregator_archive" {
  name   = format("user_data_aggregator_function-%s.zip", data.archive_file.user_data_aggregator_function.output_md5)
  bucket = google_storage_bucket.user_data_aggregator_bucket.name
  source = "${path.module}/user_data_aggregator_cloud_function.zip"
}

resource "google_service_account" "user_data_aggregator_user_service_account" {
  account_id   = "user-data-aggregator-worker"
  display_name = "User Data Aggregator Worker"
}

resource "google_service_account_iam_member" "user_data_aggregator_worker_sa" {
  service_account_id = google_service_account.user_data_aggregator_user_service_account.id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:terraform-cloud@${local.project_id}.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "user_data_aggregator_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.user_data_aggregator_role.id
  member  = format("serviceAccount:%s", google_service_account.user_data_aggregator_user_service_account.email)
}

resource "google_cloudfunctions2_function" "user_data_aggregator" {
  name        = "user-data-aggregator"
  location    = local.gcp_region
  description = format("User Data Aggregator (%s)", data.archive_file.user_data_aggregator_function.output_md5)

  build_config {
    runtime     = local.go_runtime
    entry_point = "AggregateUserData"

    source {
      storage_source {
        bucket = google_storage_bucket.user_data_aggregator_bucket.name
        object = google_storage_bucket_object.user_data_aggregator_archive.name
      }
    }
  }

  service_config {
    max_instance_count             = 1
    available_memory               = "256Mi"
    ingress_settings               = "ALLOW_INTERNAL_ONLY"
    all_traffic_on_latest_revision = true
    service_account_email          = google_service_account.user_data_aggregator_user_service_account.email

    environment_variables = {
      DINNER_DONE_BETTER_SERVICE_ENVIRONMENT = local.environment,
      # TODO: use the user_data_aggregator_user for this, currently it has permission denied for accessing tables
      # https://dba.stackexchange.com/questions/53914/permission-denied-for-relation-table
      # https://www.postgresql.org/docs/13/sql-alterdefaultprivileges.html
      DINNER_DONE_BETTER_DATABASE_USER = google_sql_user.api_user.name,
      DINNER_DONE_BETTER_DATABASE_NAME = local.database_name,
      // NOTE: if you're creating a cloud function or server for the first time, terraform cannot configure the database connection.
      // You have to go into the Cloud Run interface and deploy a new revision with a database connection, which will persist upon further deployments.
      DINNER_DONE_BETTER_DATABASE_INSTANCE_CONNECTION_NAME = google_sql_database_instance.dev.connection_name,
      GOOGLE_CLOUD_SECRET_STORE_PREFIX                     = format("projects/%d/secrets", data.google_project.project.number)
      GOOGLE_CLOUD_PROJECT_ID                              = data.google_project.project.project_id
    }

    secret_environment_variables {
      key        = "DINNER_DONE_BETTER_DATABASE_PASSWORD"
      project_id = local.project_id
      secret     = google_secret_manager_secret.api_user_database_password.secret_id
      version    = "latest"
    }

    secret_environment_variables {
      key        = "DINNER_DONE_BETTER_DATA_CHANGES_TOPIC_NAME"
      project_id = local.project_id
      secret     = google_secret_manager_secret.data_changes_topic_name.secret_id
      version    = "latest"
    }

    secret_environment_variables {
      key        = "DINNER_DONE_BETTER_USER_AGGREGATOR_TOPIC_NAME"
      project_id = local.project_id
      secret     = google_secret_manager_secret.api_user_database_password.secret_id
      version    = "latest"
    }

    secret_environment_variables {
      key        = "DINNER_DONE_BETTER_OAUTH2_TOKEN_ENCRYPTION_KEY"
      project_id = local.project_id
      secret     = google_secret_manager_secret.oauth2_token_encryption_key.secret_id
      version    = "latest"
    }
  }

  event_trigger {
    trigger_region        = local.gcp_region
    event_type            = local.pubsub_topic_publish_event
    pubsub_topic          = google_pubsub_topic.user_data_aggregator_topic.id
    retry_policy          = "RETRY_POLICY_RETRY"
    service_account_email = google_service_account.user_data_aggregator_user_service_account.email
  }
}
