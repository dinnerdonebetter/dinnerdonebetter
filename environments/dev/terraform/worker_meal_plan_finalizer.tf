resource "google_project_iam_custom_role" "meal_plan_finalizer_role" {
  role_id     = "meal_plan_finalizer_role"
  title       = "Meal Plan finalizer role"
  description = "An IAM role for the Meal Plan finalizer"
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
  ]
}

locals {
  meal_plan_finalizer_database_username = "meal_plan_finalizer_db_user"
}

resource "google_pubsub_topic" "meal_plan_finalizer_topic" {
  name = "meal_plan_finalizer_work"
}

resource "google_storage_bucket" "meal_plan_finalizer_bucket" {
  name     = "meal-plan-finalizer-cloud-function"
  location = "US"
}

data "archive_file" "meal_plan_finalizer_function" {
  type        = "zip"
  source_dir  = "${path.module}/meal_plan_finalizer_cloud_function"
  output_path = "${path.module}/meal_plan_finalizer_cloud_function.zip"
}

resource "google_storage_bucket_object" "meal_plan_finalizer_archive" {
  name   = format("meal_plan_finalizer_function-%s.zip", data.archive_file.meal_plan_finalizer_function.output_md5)
  bucket = google_storage_bucket.meal_plan_finalizer_bucket.name
  source = "${path.module}/meal_plan_finalizer_cloud_function.zip"
}

resource "google_service_account" "meal_plan_finalizer_user_service_account" {
  account_id   = "meal-plan-finalizer-worker"
  display_name = "Meal Plans Finalizer"
}

resource "google_project_iam_member" "meal_plan_finalizer_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.meal_plan_finalizer_role.id
  member  = format("serviceAccount:%s", google_service_account.meal_plan_finalizer_user_service_account.email)
}

resource "google_cloud_run_v2_job" "meal_plan_finalizer" {
  name     = "meal-plan-finalizer"
  location = local.gcp_region

  template {
    task_count  = 1
    parallelism = 1

    template {
      execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
      max_retries           = 1
      service_account       = google_service_account.meal_plan_finalizer_user_service_account.email

      volumes {
        name = "cloudsql"
        cloud_sql_instance {
          instances = [google_sql_database_instance.dev.connection_name]
        }
      }

      containers {
        image = format("%s-docker.pkg.dev/%s/%s/meal-plan-finalizer-worker", local.gcp_region, local.project_id, google_artifact_registry_repository.dev_repository.name)

        env {
          name  = "DINNER_DONE_BETTER_DATABASE_USER"
          value = google_sql_user.api_user.name
        }

        env {
          name  = "DINNER_DONE_BETTER_DATABASE_NAME"
          value = local.database_name
        }

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
          name = "DINNER_DONE_BETTER_DATABASE_PASSWORD"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.api_user_database_password.secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "DINNER_DONE_BETTER_SEGMENT_API_TOKEN"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.segment_api_token.secret_id
              version = "latest"
            }
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

resource "google_cloud_scheduler_job" "meal_plan_finalization" {
  project          = local.project_id
  region           = local.gcp_region
  name             = "meal-plan-finalizer"
  description      = "Runs the meal plan finalizer every 10 minutes"
  schedule         = "*/10 * * * *"
  time_zone        = "America/Chicago"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://${google_cloud_run_v2_job.meal_plan_finalizer.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${data.google_project.project.number}/jobs/${google_cloud_run_v2_job.meal_plan_finalizer.name}:run"

    oauth_token {
      service_account_email = google_service_account.meal_plan_finalizer_user_service_account.email
    }
  }

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_cloud_run_v2_job.meal_plan_finalizer,
  ]
}