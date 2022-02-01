resource "digitalocean_vpc" "dev" {
  name     = "prixfixe-dev-network"
  region   = local.region
  ip_range = "10.10.10.0/24"
}

resource "digitalocean_project_resources" "vpc" {
  project = digitalocean_project.prixfixe_dev.id
  resources = [
    digitalocean_vpc.dev.urn
  ]
}