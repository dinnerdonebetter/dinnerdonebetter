resource "google_project_iam_custom_role" "db_cleaner_role" {
  role_id     = "db_cleaner_role"
  title       = "database cleaner role"
  description = "An IAM role for the database cleaner"
  permissions = [
    "secretmanager.versions.access",
    "eventarc.events.receiveAuditLogWritten",
    "eventarc.events.receiveEvent",
    "run.jobs.run",
    "run.routes.invoke",
    "artifactregistry.dockerimages.get",
    "artifactregistry.dockerimages.list",
  ]
}

resource "google_service_account" "db_cleaner_user_service_account" {
  account_id   = "db-cleaner"
  display_name = "DB Cleaner"
}

resource "google_service_account_iam_member" "db_cleaner_sa" {
  service_account_id = google_service_account.outbound_emailer_user_service_account.id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:terraform-cloud@${local.project_id}.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "db_cleaner_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.db_cleaner_role.id
  member  = format("serviceAccount:%s", google_service_account.db_cleaner_user_service_account.email)
}

resource "google_cloud_run_v2_job" "db_cleaner" {
  name     = "db-cleaner"
  location = local.gcp_region

  template {
    task_count  = 1
    parallelism = 1

    template {
      execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
      max_retries           = 1
      service_account       = google_service_account.db_cleaner_user_service_account.email

      containers {
        image = format("%s-docker.pkg.dev/%s/%s/db-cleaner", local.gcp_region, local.project_id, google_artifact_registry_repository.dev_repository.name)

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi" # cannot be lower than this
          }
        }

        env {
          name  = "GOOGLE_CLOUD_SECRET_STORE_PREFIX"
          value = format("projects/%d/secrets", data.google_project.project.number)
        }

        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = data.google_project.project.project_id
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

resource "google_cloud_scheduler_job" "run_db_cleaner" {
  name             = "db-cleaner-schedule"
  description      = "Runs the database cleaner twice a month"
  schedule         = "30 1 1,15 * *"
  time_zone        = "America/Chicago"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://${google_cloud_run_v2_job.db_cleaner.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${data.google_project.project.number}/jobs/${google_cloud_run_v2_job.db_cleaner.name}:run"

    oauth_token {
      service_account_email = google_service_account.db_cleaner_user_service_account.email
    }
  }

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_cloud_run_v2_job.db_cleaner,
  ]
}