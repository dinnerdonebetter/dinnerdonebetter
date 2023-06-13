resource "google_cloudbuild_worker_pool" "pool" {
  name     = "dev"
  location = local.gcp_region
  worker_config {
    disk_size_gb   = 100
    machine_type   = "e2-medium"
    no_external_ip = false
  }
}