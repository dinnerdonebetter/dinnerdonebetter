variable "GOOGLE_CLOUD_CREDENTIALS" {}

locals {
  project_id    = "prixfixe-dev"
  gcp_region    = "us-central1"
  gcp_main_zone = "us-central1-c"
}

provider "google" {
  project     = local.project_id
  region      = local.gcp_region
  zone        = local.gcp_main_zone
  credentials = var.GOOGLE_CLOUD_CREDENTIALS
}

# you gotta enable this `cloudresourcemanager.googleapis.com` to enable the others, I think

resource "google_project_service" "iam" {
  project = local.project_id
  service = "iam.googleapis.com"
}

resource "google_project_service" "cloud_run" {
  project = local.project_id
  service = "run.googleapis.com"
}

resource "google_project_service" "container_registry" {
  project = local.project_id
  service = "containerregistry.googleapis.com"
}

resource "google_project_service" "sql_admin" {
  project = local.project_id
  service = "sqladmin.googleapis.com"
}
