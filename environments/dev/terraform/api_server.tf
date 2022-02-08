locals {
  public_url = "api.prixfixe.dev"
}

resource "google_service_account" "api_server_account" {
  account_id   = "api-server"
  display_name = "API Server"
}

resource "google_project_iam_member" "api_user" {
  project = local.project_id
  role    = "roles/viewer"
  member  = format("serviceAccount:%s", google_service_account.api_server_account.email)
}

resource "google_project_iam_binding" "api_usersecret_accessor" {
  project = local.project_id
  role    = "roles/secretmanager.secretAccessor"

  members = [
    google_project_iam_member.api_user.member,
  ]
}

resource "google_project_iam_binding" "api_usercloudsql_client" {
  project = local.project_id
  role    = "roles/cloudsql.client"

  members = [
    google_project_iam_member.api_user.member,
  ]
}

resource "google_cloud_run_service" "api_server" {
  name     = "api-server"
  location = "us-central1"

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  template {
    spec {
      service_account_name = google_service_account.api_server_account.email

      containers {
        image = "gcr.io/prixfixe-dev/api_server"

        volume_mounts {
          name       = "config"
          mount_path = "/config"
        }

        env {
          name  = "RUNNING_IN_GOOGLE_CLOUD_RUN"
          value = "true"
        }

        env {
          name  = "CONFIGURATION_FILEPATH"
          value = "/config/service-config.json"
        }

        env {
          name  = "PRIXFIXE_DATABASE_USER"
          value = local.database_username
        }

        env {
          name  = "PRIXFIXE_DATABASE_PASSWORD"
          value = random_password.database_password.result
        }

        env {
          name  = "PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME"
          value = google_sql_database_instance.dev.connection_name
        }

        env {
          name  = "PRIXFIXE_DATABASE_NAME"
          value = local.database_name
        }

        env {
          name = "PRIXFIXE_COOKIE_HASH_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.cookie_hash_key.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "PRIXFIXE_COOKIE_BLOCK_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.cookie_block_key.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "PRIXFIXE_PASETO_LOCAL_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.paseto_local_key.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "PRIXFIXE_DATA_CHANGES_TOPIC"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.data_changes_topic_name.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "PRIXFIXE_SENDGRID_API_TOKEN"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.sendgrid_api_token.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "PRIXFIXE_SEGMENT_API_TOKEN"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.segment_api_token.secret_id
              key  = "latest"
            }
          }
        }

      }

      volumes {
        name = "config"
        secret {
          secret_name  = google_secret_manager_secret.api_service_config.secret_id
          default_mode = 256 # 0400
          items {
            key  = "latest"
            path = "service-config.json"
            mode = 256 # 0400
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "1"
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.dev.connection_name
        "run.googleapis.com/client-name"        = "terraform"
      }
    }
  }
}
