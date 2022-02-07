variable "GOOGLE_CLOUD_CREDENTIALS" {}

provider "google" {
  project     = "prixfixe-dev"
  region      = "us-central1"
  zone        = "us-central1-c"
  credentials = var.GOOGLE_CLOUD_CREDENTIALS
}
