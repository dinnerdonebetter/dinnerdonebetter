variable "GCP_CREDENTIALS_DEV" {}

provider "google" {
  project     = "prixfixe-dev"
  region      = "us-central1"
  zone        = "us-central1-c"
  credentials = var.GCP_CREDENTIALS_DEV
}

resource "google_container_registry" "registry" {
  location = "US"
}