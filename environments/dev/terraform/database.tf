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

resource "google_sql_database" "api_database" {
  name     = local.database_name
  instance = google_sql_database_instance.dev.name
}