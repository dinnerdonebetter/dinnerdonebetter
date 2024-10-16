resource "google_artifact_registry_repository" "dev_repository" {
  location      = local.gcp_region
  repository_id = "containers"
  description   = "the container images for the dev environment"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "db-cleaner-container" {
  location      = local.gcp_region
  repository_id = "db-cleaner"
  description   = "database cleaner worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "search-index-scheduler-container" {
  location      = local.gcp_region
  repository_id = "search-index-scheduler"
  description   = "search data index scheduler worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "email-prober-container" {
  location      = local.gcp_region
  repository_id = "email-prober"
  description   = "email prober worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "meal-plan-finalizer-container" {
  location      = local.gcp_region
  repository_id = "meal-plan-finalizer"
  description   = "meal plan finalizer worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "meal-plan-grocery-list-initializer-container" {
  location      = local.gcp_region
  repository_id = "meal-plan-grocery-list-initializer"
  description   = "mela plan grocery list initializer worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "meal-plan-task-creator-container" {
  location      = local.gcp_region
  repository_id = "meal-plan-task-creator"
  description   = "meal plan task creator worker image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "data-changes-handler-container" {
  location      = local.gcp_region
  repository_id = "data-changes-handler"
  description   = "data changes handler function image"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "email-handler-container" {
  location      = local.gcp_region
  repository_id = "email-handler"
  description   = "email handler function image"
  format        = "DOCKER"
}


resource "google_artifact_registry_repository" "search-indexer-container" {
  location      = local.gcp_region
  repository_id = "search-indexer"
  description   = "search indexing function image"
  format        = "DOCKER"
}


resource "google_artifact_registry_repository" "webhook-executor-container" {
  location      = local.gcp_region
  repository_id = "webhook-executor"
  description   = "webhook execution function image"
  format        = "DOCKER"
}
