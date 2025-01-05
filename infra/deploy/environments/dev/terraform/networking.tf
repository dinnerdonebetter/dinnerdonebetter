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

resource "google_certificate_manager_dns_authorization" "api" {
  name        = "api"
  domain      = "api.dinnerdonebetter.dev"
}

resource "google_certificate_manager_dns_authorization" "admin" {
  name        = "admin"
  domain      = "admin.dinnerdonebetter.dev"
}

resource "google_certificate_manager_dns_authorization" "app" {
  name        = "app"
  domain      = "app.dinnerdonebetter.dev"
}

resource "google_certificate_manager_dns_authorization" "www" {
  name   = "www"
  domain = "www.dinnerdonebetter.dev"
}

resource "google_certificate_manager_dns_authorization" "media" {
  name   = "media"
  domain = "media.dinnerdonebetter.dev"
}

resource "google_certificate_manager_dns_authorization" "userdata" {
  name   = "userdata"
  domain = "userdata.dinnerdonebetter.dev"
}

resource "google_certificate_manager_certificate" "default" {
  name        = "dns-cert"
  description = "The default cert"
  scope       = "EDGE_CACHE"
  labels = {
    (local.managed_by_label) = "terraform"
  }
  managed {
    domains = [
      google_certificate_manager_dns_authorization.api.domain,
      google_certificate_manager_dns_authorization.admin.domain,
      google_certificate_manager_dns_authorization.app.domain,
      google_certificate_manager_dns_authorization.www.domain,
      google_certificate_manager_dns_authorization.media.domain,
      google_certificate_manager_dns_authorization.userdata.domain,
    ]
    dns_authorizations = [
      google_certificate_manager_dns_authorization.api.id,
      google_certificate_manager_dns_authorization.admin.id,
      google_certificate_manager_dns_authorization.app.id,
      google_certificate_manager_dns_authorization.www.id,
      google_certificate_manager_dns_authorization.media.id,
      google_certificate_manager_dns_authorization.userdata.id,
    ]
  }
}

resource "cloudflare_record" "domain_verification_records" {
  for_each = {
    "0": google_certificate_manager_dns_authorization.api
    "1": google_certificate_manager_dns_authorization.admin
    "2": google_certificate_manager_dns_authorization.app
    "3": google_certificate_manager_dns_authorization.www
    "4": google_certificate_manager_dns_authorization.media
    "5": google_certificate_manager_dns_authorization.userdata
  }

  zone_id  = var.CLOUDFLARE_ZONE_ID
  name     = each.value.name
  type     = upper(each.value.type)
  content  = each.value.data
  ttl      = 1
  proxied  = false
  comment  = "Managed by Terraform"
}
