# Segment API token

variable "SEGMENT_API_TOKEN" {}
resource "google_secret_manager_secret" "segment_api_token" {
  secret_id = "webapp_segment_api_token"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "segment_api_token" {
  secret = google_secret_manager_secret.segment_api_token.id

  secret_data = var.SEGMENT_API_TOKEN
}

# Service OAuth2 Client ID

variable "DINNER_DONE_BETTER_OAUTH2_CLIENT_ID" {}
resource "google_secret_manager_secret" "ddb_oauth2_client_id" {
  secret_id = "webapp_ddb_oauth2_client_id"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "ddb_oauth2_client_id" {
  secret = google_secret_manager_secret.ddb_oauth2_client_id.id

  secret_data = var.DINNER_DONE_BETTER_OAUTH2_CLIENT_ID
}

# Service OAuth2 Client Secret

variable "DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET" {}
resource "google_secret_manager_secret" "ddb_oauth2_client_secret" {
  secret_id = "webapp_ddb_oauth2_client_secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "ddb_oauth2_client_secret" {
  secret = google_secret_manager_secret.ddb_oauth2_client_secret.id

  secret_data = var.DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET
}

# Cookie encryption key

resource "random_string" "cookie_encryption_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "cookie_encryption_key" {
  secret_id = "webapp_cookie_encryption_key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "cookie_encryption_key" {
  secret = google_secret_manager_secret.cookie_encryption_key.id

  secret_data = random_string.cookie_encryption_key.result
}

# Cookie encryption initialization vector

resource "random_bytes" "cookie_encryption_iv" {
  length = 32
}

resource "google_secret_manager_secret" "cookie_encryption_iv" {
  secret_id = "webapp_cookie_encryption_iv"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "cookie_encryption_iv" {
  secret = google_secret_manager_secret.cookie_encryption_iv.id

  secret_data = random_bytes.cookie_encryption_iv.base64
}
