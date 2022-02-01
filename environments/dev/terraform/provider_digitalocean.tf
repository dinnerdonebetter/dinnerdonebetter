variable "DIGITAL_OCEAN_API_TOKEN" {}
variable "SPACES_ACCESS_KEY_ID" {}
variable "SPACES_SECRET_ACCESS_KEY" {}

locals {
  region = "nyc3"
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.DIGITAL_OCEAN_API_TOKEN

  spaces_access_id  = var.SPACES_ACCESS_KEY_ID
  spaces_secret_key = var.SPACES_SECRET_ACCESS_KEY
}

resource "digitalocean_container_registry" "dev" {
  name                   = "dev-prixfixe"
  subscription_tier_slug = "starter"
}

output "docr_registry_domain" {
  value = digitalocean_container_registry.dev.endpoint
}

resource "digitalocean_project" "prixfixe_dev" {
  name        = "prixfixe-dev"
  description = "the dev environment for PrixFixe"
  purpose     = "Service or API"
  environment = "Development"
  resources = [
    digitalocean_kubernetes_cluster.dev.id,
    digitalocean_database_cluster.database.urn,
    digitalocean_spaces_bucket.config.urn,
  ]
}