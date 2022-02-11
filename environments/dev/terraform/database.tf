locals {
  database_name = "prixfixe"
}

resource "google_sql_ssl_cert" "client_cert" {
  common_name = "prixfixe-dev"
  instance    = google_sql_database_instance.dev.name
}

resource "google_sql_database_instance" "dev" {
  name                = "dev-whatever-1644302117"
  database_version    = "POSTGRES_13"
  region              = local.gcp_region
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "api_database" {
  name     = local.database_name
  instance = google_sql_database_instance.dev.name
}
