resource "digitalocean_vpc" "dev" {
  name     = "prixfixe-dev-network"
  region   = local.region
  ip_range = "10.10.10.0/24"
}
