locals {
  database_name = "dinner-done-better"
}

# BEGIN Networking shit

resource "google_compute_network" "private_network" {
  provider = google

  name = "private-network"
}

resource "google_compute_global_address" "private_ip_address" {
  provider = google

  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider = google

  network                 = google_compute_network.private_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}


# END Networking shit

resource "google_sql_database_instance" "dev" {
  name                = "dev"
  database_version    = "POSTGRES_17"
  region              = local.gcp_region
  deletion_protection = false

  depends_on = [google_service_networking_connection.private_vpc_connection]

  settings {
    tier                  = "db-f1-micro"
    disk_size             = 20
    disk_autoresize       = true
    disk_autoresize_limit = 50
    edition               = "ENTERPRISE"

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 4096
      record_application_tags = false
      record_client_address   = false
    }

    maintenance_window {
      day          = 7
      hour         = 3
      update_track = "stable"
    }

    ip_configuration {
      ssl_mode                                      = "ENCRYPTED_ONLY"
      ipv4_enabled                                  = false
      private_network                               = google_compute_network.private_network.self_link
      enable_private_path_for_google_cloud_services = true
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
