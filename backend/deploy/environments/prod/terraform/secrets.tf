# API server oauth2 token encryption key
resource "random_string" "oauth2_token_encryption_key" {
  length  = 32
  special = false
}

# API server user device token encryption key
resource "random_string" "user_device_token_encryption_key" {
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

# Consumer webapp cookie encryption keys (base64-encoded for COOKIE_HASH_KEY / COOKIE_BLOCK_KEY)
resource "random_string" "consumer_webapp_cookie_hash_key" {
  length  = 32
  special = false
}
resource "random_string" "consumer_webapp_cookie_block_key" {
  length  = 32
  special = false
}

### External API services ###

# Sendgrid token
variable "SENDGRID_API_KEY" {}

# Resend API key
variable "RESEND_API_KEY" {}

# Segment API token
variable "SEGMENT_API_TOKEN" {}

# PostHog API tokens: Project API Key (events); Personal API Key (feature flags API).
# Add both to Terraform Cloud: Workspace → Variables → Add variable.
variable "POSTHOG_API_KEY" {}
variable "POSTHOG_PERSONAL_API_KEY" {}

# Grafana Cloud API keys
# NOTE: the passwords are all effectively the same, but they maybe won't be one day? Who knows.
variable "GRAFANA_CLOUD_PROMETHEUS_USERNAME" {}
variable "GRAFANA_CLOUD_PROMETHEUS_PASSWORD" {}
variable "GRAFANA_CLOUD_LOKI_USERNAME" {}
variable "GRAFANA_CLOUD_LOKI_PASSWORD" {}
variable "GRAFANA_CLOUD_TEMPO_USERNAME" {}
variable "GRAFANA_CLOUD_TEMPO_PASSWORD" {}
variable "GRAFANA_CLOUD_PYROSCOPE_USERNAME" {}
variable "GRAFANA_CLOUD_PYROSCOPE_PASSWORD" {}

variable "APNS_KEY_ID" {}
variable "APNS_AUTH_KEY_P8" {
  sensitive = true
}
variable "APNS_TEAM_ID" {
  default = "K8R2Q5UWQS"
}
variable "APNS_BUNDLE_ID" {
  default = "com.dinnerdonebetter.ios"
}
variable "APNS_PRODUCTION" {
  default     = "false"
  description = "Use APNs production environment (true) or sandbox (false). Sandbox for debug/TestFlight builds; production for App Store."
}

# Admin webapp config (cookie name and domain - required for admin webapp)
variable "ADMIN_WEBAPP_COOKIE_NAME" {
  default = "admin_webapp"
}
variable "ADMIN_WEBAPP_COOKIE_DOMAIN" {
  default = "admin.dinnerdonebetter.com"
}

# Consumer webapp config (cookie name and domain - required for consumer webapp / root site)
variable "CONSUMER_WEBAPP_COOKIE_NAME" {
  default = "consumer_session"
}
variable "CONSUMER_WEBAPP_COOKIE_DOMAIN" {
  default = "dinnerdonebetter.com"
}
