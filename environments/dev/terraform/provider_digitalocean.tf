variable "DIGITAL_OCEAN_API_TOKEN" {}
#variable "SPACES_ACCESS_KEY_ID" {}
#variable "SPACES_SECRET_ACCESS_KEY" {}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.DIGITAL_OCEAN_API_TOKEN
}

resource "digitalocean_container_registry" "dev" {
  name                   = "dev"
  subscription_tier_slug = "starter"
}

output "docr_registry_domain" {
  value = digitalocean_container_registry.dev.endpoint
}