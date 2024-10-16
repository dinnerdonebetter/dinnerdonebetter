resource "google_artifact_registry_repository" "dev_repository" {
  location      = local.gcp_region
  repository_id = "containers"
  description   = "the container images for the dev environment"
  format        = "DOCKER"

  cleanup_policies {
    id     = "keep-5-latests"
    action = "KEEP"
    most_recent_versions {
      package_name_prefixes = ["latest"]
      keep_count            = 5
    }
  }
}
