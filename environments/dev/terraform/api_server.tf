locals {
  api_database_username = "api_db_user"
  public_url            = "api.prixfixe.dev"
}

resource "google_project_iam_custom_role" "api_server_role" {
  role_id     = "api_server_role"
  title       = "API Server role"
  description = "An IAM role for the API server"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
  ]
}

resource "google_service_account" "api_user_service_account" {
  account_id   = "api-server"
  display_name = "API Server"
}

resource "google_project_iam_member" "api_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.api_server_role.id
  member  = format("serviceAccount:%s", google_service_account.api_user_service_account.email)
}

resource "google_project_iam_binding" "api_user_secret_accessor" {
  project = local.project_id
  role    = "roles/iam.serviceAccountUser"

  members = [
    google_project_iam_member.api_user.member,
  ]
}

# this allows the service to be on the public internet
data "google_iam_policy" "public_access" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "public_access" {
  location = google_cloud_run_service.api_server.location
  project  = google_cloud_run_service.api_server.project
  service  = google_cloud_run_service.api_server.name

  policy_data = data.google_iam_policy.public_access.policy_data
}

resource "random_password" "api_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "api_user" {
  name     = local.api_database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.api_user_database_password.result
}

resource "google_cloud_run_service" "api_server" {
  name     = "api-server"
  location = local.gcp_region

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true

  template {
    spec {
      service_account_name = google_service_account.api_user_service_account.email

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
          name  = "GOOGLE_CLOUD_SECRET_STORE_PREFIX"
          value = format("projects/%d/secrets", data.google_project.project.number)
        }

        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = data.google_project.project.project_id
        }

        env {
          name  = "CONFIGURATION_FILEPATH"
          value = "/config/service-config.json"
        }

        env {
          name  = "PRIXFIXE_DATABASE_USER"
          value = local.api_database_username
        }

        env {
          name  = "PRIXFIXE_DATABASE_PASSWORD"
          value = random_password.api_user_database_password.result
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

resource "cloudflare_record" "api_cname_record" {
  zone_id = var.CLOUDFLARE_ZONE_ID
  name    = "api.prixfixe.dev"
  type    = "CNAME"
  value   = "ghs.googlehosted.com."
  ttl     = 1
  proxied = true
}
