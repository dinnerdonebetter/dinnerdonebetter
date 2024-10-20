
resource "google_storage_bucket_access_control" "public_rule" {
  bucket = google_storage_bucket.api_media.name
  role   = "READER"
  entity = "allUsers"
}


data "google_iam_policy" "public_policy" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
    ]
  }
}
