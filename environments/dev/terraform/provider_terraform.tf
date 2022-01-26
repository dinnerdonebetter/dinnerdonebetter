variable "GOOGLE_CREDENTIALS" {}

provider "google" {
  project = "prixfixe-dev"
  region  = "us-central1"
  zone    = "us-central1-c"
}