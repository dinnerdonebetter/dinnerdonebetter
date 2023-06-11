resource "google_project_iam_custom_role" "search_data_index_scheduler_role" {
  role_id     = "search_data_index_scheduler_role"
  title       = "Search data index scheduler role"
  description = "An IAM role for the Search data index scheduler"
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
  search_data_index_scheduler_database_username = "search_data_index_scheduler_db_user"
}

#resource "google_pubsub_topic" "search_data_index_scheduler_topic" {
#  name = "search_data_index_scheduler_work"
#}

#resource "google_cloud_scheduler_job" "search_data_index_scheduling" {
#  project = local.project_id
#  region  = local.gcp_region
#  name    = "search-data-index-scheduler"
#
#  schedule  = "*/10 * * * *" # every ten minutes
#  time_zone = "America/Chicago"
#
#  pubsub_target {
#    topic_name = google_pubsub_topic.search_data_index_scheduler_topic.id
#    data       = base64encode("{}")
#  }
#}

resource "google_storage_bucket" "search_data_index_scheduler_bucket" {
  name     = "search-data-index-scheduler-cloud-function"
  location = "US"
}

data "archive_file" "search_data_index_scheduler_function" {
  type        = "zip"
  source_dir  = "${path.module}/search_data_index_scheduler_cloud_function"
  output_path = "${path.module}/search_data_index_scheduler_cloud_function.zip"
}

resource "google_storage_bucket_object" "search_data_index_scheduler_archive" {
  name   = format("search_data_index_scheduler_function-%s.zip", data.archive_file.search_data_index_scheduler_function.output_md5)
  bucket = google_storage_bucket.search_data_index_scheduler_bucket.name
  source = "${path.module}/search_data_index_scheduler_cloud_function.zip"
}

resource "google_service_account" "search_data_index_scheduler_user_service_account" {
  account_id   = "search-data-index-scheduler"
  display_name = "Search Data Index Scheduler"
}

resource "google_project_iam_member" "search_data_index_scheduler_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.search_data_index_scheduler_role.id
  member  = format("serviceAccount:%s", google_service_account.search_data_index_scheduler_user_service_account.email)
}

resource "random_password" "search_data_index_scheduler_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_secret_manager_secret" "search_data_index_scheduler_user_database_password" {
  secret_id = "search_data_index_scheduler_user_database_password"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "search_data_index_scheduler_user_database_password" {
  secret = google_secret_manager_secret.search_data_index_scheduler_user_database_password.id

  secret_data = random_password.search_data_index_scheduler_user_database_password.result
}

resource "google_sql_user" "search_data_index_scheduler_user" {
  name     = local.search_data_index_scheduler_database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.search_data_index_scheduler_user_database_password.result
}

# Permissions on the service account used by the function and Eventarc trigger
resource "google_project_iam_member" "search_data_index_scheduler_invoking" {
  project = local.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.search_data_index_scheduler_user_service_account.email}"
}

resource "google_project_iam_member" "search_data_index_scheduler_event_receiving" {
  project    = local.project_id
  role       = "roles/eventarc.eventReceiver"
  member     = "serviceAccount:${google_service_account.search_data_index_scheduler_user_service_account.email}"
  depends_on = [google_project_iam_member.search_data_index_scheduler_invoking]
}

resource "google_project_iam_member" "search_data_index_scheduler_artifactregistry_reader" {
  project    = local.project_id
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${google_service_account.search_data_index_scheduler_user_service_account.email}"
  depends_on = [google_project_iam_member.search_data_index_scheduler_event_receiving]
}

#resource "google_cloudfunctions2_function" "search_data_index_scheduler" {
#  depends_on = [
#    google_cloud_scheduler_job.search_data_index_scheduling,
#    google_project_iam_member.search_data_index_scheduler_event_receiving,
#    google_project_iam_member.search_data_index_scheduler_artifactregistry_reader,
#  ]
#
#  name        = "search-data-index-scheduler"
#  description = "Search Data Index Scheduler"
#  location    = local.gcp_region
#
#  build_config {
#    runtime     = local.go_runtime
#    entry_point = "ScheduleIndexOperation"
#
#    source {
#      storage_source {
#        bucket = google_storage_bucket.search_data_index_scheduler_bucket.name
#        object = google_storage_bucket_object.search_data_index_scheduler_archive.name
#      }
#    }
#  }
#
#  service_config {
#    available_memory               = "128Mi"
#    ingress_settings               = "ALLOW_INTERNAL_ONLY"
#    all_traffic_on_latest_revision = true
#    service_account_email          = google_service_account.search_data_index_scheduler_user_service_account.email
#
#    environment_variables = {
#      DINNER_DONE_BETTER_SERVICE_ENVIRONMENT = local.environment,
#      # TODO: use the search_data_index_scheduler_user for this, currently it has permission denied for accessing tables
#      # https://dba.stackexchange.com/questions/53914/permission-denied-for-relation-table
#      # https://www.postgresql.org/docs/13/sql-alterdefaultprivileges.html
#      DINNER_DONE_BETTER_DATABASE_USER = google_sql_user.api_user.name,
#      DINNER_DONE_BETTER_DATABASE_NAME = local.database_name,
#      // NOTE: if you're creating a cloud function or server for the first time, terraform cannot configure the database connection.
#      // You have to go into the Cloud Run interface and deploy a new revision with a database connection, which will persist upon further deployments.
#      DINNER_DONE_BETTER_DATABASE_INSTANCE_CONNECTION_NAME = google_sql_database_instance.dev.connection_name,
#      GOOGLE_CLOUD_SECRET_STORE_PREFIX                     = format("projects/%d/secrets", data.google_project.project.number)
#      GOOGLE_CLOUD_PROJECT_ID                              = data.google_project.project.project_id
#      DATA_CHANGES_TOPIC_NAME                              = google_pubsub_topic.data_changes_topic.name
#      SEARCH_INDEXING_TOPIC_NAME                           = google_pubsub_topic.search_index_requests_topic.name
#    }
#
#    secret_environment_variables {
#      key        = "DINNER_DONE_BETTER_DATABASE_PASSWORD"
#      project_id = local.project_id
#      secret     = google_secret_manager_secret.api_user_database_password.secret_id
#      version    = "latest"
#    }
#  }
#
#  event_trigger {
#    trigger_region        = local.gcp_region
#    event_type            = local.pubsub_topic_publish_event
#    pubsub_topic          = google_pubsub_topic.search_data_index_scheduler_topic.id
#    retry_policy          = "RETRY_POLICY_RETRY"
#    service_account_email = google_service_account.search_data_index_scheduler_user_service_account.email
#  }
#}

resource "google_artifact_registry_repository" "index_scheduler_repository" {
  location      = local.gcp_region
  repository_id = "search-data-index-scheduler"
  description   = "the container image for the search data index scheduler"
  format        = "DOCKER"

  docker_config {
    immutable_tags = true
  }
}

resource "google_cloud_run_v2_job" "search_data_index_scheduler" {
  name     = "search-data-index-scheduler"
  location = local.gcp_region

  template {
    task_count  = 1
    parallelism = 1

    template {
      execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
      max_retries           = 1

      volumes {
        name = "cloudsql"
        cloud_sql_instance {
          instances = [google_sql_database_instance.dev.connection_name]
        }
      }

      containers {
        image = format("gcr.io/dinner-done-better-dev/%s", google_artifact_registry_repository.index_scheduler_repository.name)

        env {
          name  = "DINNER_DONE_BETTER_DATABASE_INSTANCE_CONNECTION_NAME"
          value = google_sql_database_instance.dev.connection_name
        }

        env {
          name  = "GOOGLE_CLOUD_SECRET_STORE_PREFIX"
          value = format("projects/%d/secrets", data.google_project.project.number)
        }

        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = data.google_project.project.project_id
        }

        env {
          name  = "DATA_CHANGES_TOPIC_NAME"
          value = google_pubsub_topic.data_changes_topic.name
        }

        env {
          name  = "SEARCH_INDEXING_TOPIC_NAME"
          value = google_pubsub_topic.search_index_requests_topic.name
        }

        env {
          name = "DINNER_DONE_BETTER_DATABASE_PASSWORD"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.api_user_database_password.secret_id
              version = "latest"
            }
          }
        }

        resources {
          limits {
            cpu    = "1"
            memory = "128Mi"
          }
        }

        volume_mounts {
          name       = "cloudsql"
          mount_path = "/cloudsql"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      launch_stage,
    ]
  }
}

resource "google_cloud_scheduler_job" "run_data_index_scheduler" {
  name             = "scheduled-data-indexing"
  description      = "Runs the search data index scheduler every 10 minutes"
  schedule         = "*/10 * * * *"
  time_zone        = "America/Chicago"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }


  http_target {
    http_method = "POST"
    uri         = "https://${google_cloud_run_v2_job.search_data_index_scheduler.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${data.google_project.project.number}/jobs/${google_cloud_run_v2_job.search_data_index_scheduler.name}:run"

    oauth_token {
      service_account_email = google_service_account.search_data_index_scheduler_user_service_account.email
    }
  }


  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_cloud_run_v2_job.search_data_index_scheduler,
  ]
}