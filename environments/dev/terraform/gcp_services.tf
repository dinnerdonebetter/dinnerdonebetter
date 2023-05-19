resource "google_project_service" "cloud_run_api" {
  service = "run.googleapis.com"
}
resource "google_project_service" "container_registry" {
  service = "containerregistry.googleapis.com"
}
