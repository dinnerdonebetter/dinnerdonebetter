# API server oauth2 token encryption key
resource "random_string" "oauth2_token_encryption_key" {
  length  = 32
  special = false
}

# JWT Signing key
resource "random_string" "jwt_signing_key" {
  length  = 32
  special = false
}

# Admin webapp cookie encryption keys (base64-encoded for COOKIE_HASH_KEY / COOKIE_BLOCK_KEY)
resource "random_string" "admin_webapp_cookie_hash_key" {
  length  = 32
  special = false
}
resource "random_string" "admin_webapp_cookie_block_key" {
  length  = 32
  special = false
}

### External API services ###

# Sendgrid token
variable "SENDGRID_API_KEY" {}

# Segment API token
variable "SEGMENT_API_TOKEN" {}

# PostHog API tokens (Project API Key for events; Personal API Key for feature flags API)
variable "POSTHOG_API_KEY" {}

# Grafana Cloud API keys
# NOTE: the passwords are all effectively the same, but they maybe won't be one day? Who knows.
variable "GRAFANA_CLOUD_PROMETHEUS_USERNAME" {}
variable "GRAFANA_CLOUD_PROMETHEUS_PASSWORD" {}
variable "GRAFANA_CLOUD_LOKI_USERNAME" {}
variable "GRAFANA_CLOUD_LOKI_PASSWORD" {}
variable "GRAFANA_CLOUD_TEMPO_USERNAME" {}
variable "GRAFANA_CLOUD_TEMPO_PASSWORD" {}

# Admin webapp config (cookie name and domain - required for admin webapp)
variable "ADMIN_WEBAPP_COOKIE_NAME" {
  default = "admin_webapp"
}
variable "ADMIN_WEBAPP_COOKIE_DOMAIN" {
  default = "admin.dinnerdonebetter.com"
}
