locals {
  database_username = "prixfixe_api"
  database_name     = "prixfixe"
}

resource "random_password" "database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_database_instance" "dev" {
  name             = "dev"
  database_version = "POSTGRES_13"
  region           = "us-central1"

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
  host     = "prixfixe.dev"
  password = random_password.database_password.result
}