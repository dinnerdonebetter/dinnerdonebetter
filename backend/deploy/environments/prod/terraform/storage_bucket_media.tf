resource "google_storage_bucket" "api_media" {
  provider                    = google
  name                        = "dinner-done-better-prod-media"
  location                    = "US"
  uniform_bucket_level_access = false
  force_destroy               = true

  versioning {
    enabled = true
  }

  website {
    main_page_suffix = "index.html"
    not_found_page   = "index.html"
  }

  cors {
    origin          = ["https://dinnerdonebetter.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age                = 7
      num_newer_versions = 3
    }
  }
}

resource "google_storage_bucket_iam_policy" "api_media_policy" {
  bucket      = google_storage_bucket.api_media.name
  policy_data = data.google_iam_policy.api_media_policy.policy_data
}

# NOTE: media.dinnerdonebetter.com CNAME removed. Bucket uses storage.googleapis.com URLs.
# To serve at media.dinnerdonebetter.com, verify domain ownership in Search Console first,
# then create a bucket named media.dinnerdonebetter.com and restore this record.
