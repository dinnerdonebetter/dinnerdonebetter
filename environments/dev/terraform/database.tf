locals {
  database_username = "api_db_user"
  database_name     = "prixfixe"
}

resource "random_password" "database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_database_instance" "dev" {
  name                = "dev-whatever-1644287446"
  database_version    = "POSTGRES_13"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_ssl_cert" "client_cert" {
  common_name = "prixfixe-dev"
  instance    = google_sql_database_instance.dev.name
}

resource "google_sql_user" "users" {
  name     = local.database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.database_password.result
}

# Database connection string

resource "google_secret_manager_secret" "database_connection_string" {
  secret_id = "database_connection_string"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "database_connection_string" {
  secret = google_secret_manager_secret.database_connection_string.id

  secret_data = format(
    "user=%s dbname=%s password='%s' host=%s port=5432",
    local.database_username,
    local.database_name,
    random_password.database_password.result,
    google_sql_database_instance.dev.ip_address.0.ip_address,
  )
}