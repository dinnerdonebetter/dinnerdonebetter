resource "google_artifact_registry_repository" "dev_repository" {
  location      = local.gcp_region
  repository_id = "containers"
  description   = "the container images for the dev environment"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "db-cleaner-container" {
  location      = "us-central1"
  repository_id = "db-cleaner"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "search-index-scheduler-container" {
  location      = "us-central1"
  repository_id = "search-index-scheduler"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "email-prober-container" {
  location      = "us-central1"
  repository_id = "email-prober"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "meal-plan-finalizer-container" {
  location      = "us-central1"
  repository_id = "meal-plan-finalizer"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "meal-plan-grocery-list-initializer-container" {
  location      = "us-central1"
  repository_id = "meal-plan-grocery-list-initializer"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "meal-plan-task-creator-container" {
  location      = "us-central1"
  repository_id = "meal-plan-task-creator"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}
