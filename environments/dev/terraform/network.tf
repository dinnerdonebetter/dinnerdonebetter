resource "digitalocean_vpc" "dev" {
  name     = "prixfixe-dev-network"
  region   = local.region
  ip_range = "10.10.10.0/24"
}

resource "digitalocean_loadbalancer" "public" {
  name                  = "loadbalancer-1"
  region                = local.region
  size                  = "lb-small"
  size_unit             = 1
  enable_proxy_protocol = true

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 8000
    target_protocol = "http"
  }

  droplet_ids = []
  vpc_uuid    = digitalocean_vpc.dev.id
}