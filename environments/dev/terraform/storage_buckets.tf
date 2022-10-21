resource "google_storage_bucket" "image_uploads" {
  name          = "media.prixfixe.dev"
  location      = "US"
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["https://media.prixfixe.dev"]
    method          = ["GET"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
}