variable "GOOGLE_CLOUD_CREDENTIALS" {}

locals {
  project_id    = "dinner-done-better-dev"
  gcp_region    = "us-central1"
  gcp_main_zone = "us-central1-b"
}

provider "google" {
  project     = local.project_id
  region      = local.gcp_region
  zone        = local.gcp_main_zone
  credentials = var.GOOGLE_CLOUD_CREDENTIALS
}

data "google_project" "project" {}
