locals {
  public_url = "api.prixfixe.dev"
}

resource "google_container_registry" "registry" {
  location = "US"
}