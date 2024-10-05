# Cookie encryption key

resource "random_string" "cookie_encryption_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "cookie_encryption_key" {
  secret_id = "admin_cookie_encryption_key"

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
  secret_id = "admin_cookie_encryption_iv"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "cookie_encryption_iv" {
  secret = google_secret_manager_secret.cookie_encryption_iv.id

  secret_data = random_bytes.cookie_encryption_iv.base64
}
