# Segment API token

variable "SEGMENT_API_TOKEN" {}
resource "google_secret_manager_secret" "segment_api_token" {
  secret_id = "landing_segment_api_token"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "segment_api_token" {
  secret = google_secret_manager_secret.segment_api_token.id

  secret_data = var.SEGMENT_API_TOKEN
}
