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

resource "google_compute_address" "static_ip" {
  name = "dev"
  labels = {
    (local.managed_by_label) = "terraform"
  }
}

resource "google_compute_managed_ssl_certificate" "dev" {
  name = "dev-cert"

  managed {
    domains = ["dinnerdonebetter.dev"]
  }
}
