variable "GOOGLE_CLOUD_CREDENTIALS" {}

locals {
  project_id = "prixfixe-dev"
}

provider "google" {
  project     = local.project_id
  region      = "us-central1"
  zone        = "us-central1-c"
  credentials = var.GOOGLE_CLOUD_CREDENTIALS
}

resource "google_project_service" "iam" {
  project = local.project_id
  service = "iam.googleapis.com"
}

resource "google_project_service" "run" {
  project = local.project_id
  service = "run.googleapis.com"
}

resource "google_project_service" "containerregistry" {
  project = local.project_id
  service = "containerregistry.googleapis.com"
}


resource "google_project_service" "sqladmin" {
  project = local.project_id
  service = "sqladmin.googleapis.com"
}
