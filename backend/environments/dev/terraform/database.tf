locals {
  database_name = "dinner-done-better"
}

resource "google_sql_database_instance" "dev" {
  name                = "dev"
  database_version    = "POSTGRES_17"
  region              = local.gcp_region
  deletion_protection = false

  settings {
    tier                  = "db-f1-micro"
    disk_size             = 20
    disk_autoresize       = true
    disk_autoresize_limit = 50

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 4096
      record_application_tags = false
      record_client_address   = false
    }

    ip_configuration {
      ssl_mode                                      = "ENCRYPTED_ONLY"
      enable_private_path_for_google_cloud_services = true
    }

    maintenance_window {
      hour         = 3 # 2AM
      day          = 7 # Sunday
      update_track = "stable"
    }
  }
}

resource "google_sql_ssl_cert" "client_cert" {
  common_name = "dinner-done-better-dev"
  instance    = google_sql_database_instance.dev.name
}

resource "google_sql_database" "api_database" {
  name     = local.database_name
  instance = google_sql_database_instance.dev.name
}
