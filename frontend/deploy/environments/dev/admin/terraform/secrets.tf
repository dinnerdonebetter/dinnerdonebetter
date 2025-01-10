# Cookie encryption key
resource "random_string" "cookie_encryption_key" {
  length  = 32
  special = false
}

# Cookie encryption initialization vector
resource "random_bytes" "cookie_encryption_iv" {
  length = 32
}

# Service OAuth2 Client ID & Secret
variable "DINNER_DONE_BETTER_OAUTH2_CLIENT_ID" {}
variable "DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET" {}

# Segment API token
variable "SEGMENT_API_TOKEN" {}
