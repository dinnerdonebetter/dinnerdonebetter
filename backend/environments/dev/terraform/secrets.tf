# API config

resource "google_secret_manager_secret" "api_service_config" {
  secret_id = "api_service_config"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "api_service_config" {
  secret = google_secret_manager_secret.api_service_config.id

  secret_data = file("${path.module}/service-config.json")
}

# Data changes topic

resource "google_secret_manager_secret" "data_changes_topic_name" {
  secret_id = "data_changes_topic_name"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "data_changes_topic_name" {
  secret = google_secret_manager_secret.data_changes_topic_name.id

  secret_data = google_pubsub_topic.data_changes_topic.name
}

# outbound emails

resource "google_secret_manager_secret" "outbound_emails_topic_name" {
  secret_id = "outbound_emails_topic_name"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "outbound_emails_topic_name" {
  secret = google_secret_manager_secret.outbound_emails_topic_name.id

  secret_data = google_pubsub_topic.outbound_emails_topic.name
}

# data aggregation

resource "google_secret_manager_secret" "data_aggregation_topic_name" {
  secret_id = "data_aggregation_topic_name"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "data_aggregation_topic_names" {
  secret = google_secret_manager_secret.data_aggregation_topic_name.id

  secret_data = google_pubsub_topic.user_data_aggregator_topic.name
}

# API server oauth2 token encryption key

resource "random_string" "oauth2_token_encryption_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "oauth2_token_encryption_key" {
  secret_id = "oauth2_token_encryption_key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "oauth2_token_encryption_key" {
  secret = google_secret_manager_secret.oauth2_token_encryption_key.id

  secret_data = random_string.oauth2_token_encryption_key.result
}

# JWT Signing key

resource "random_string" "jwt_signing_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "jwt_signing_key" {
  secret_id = "jwt_signing_key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "jwt_signing_key" {
  secret = google_secret_manager_secret.jwt_signing_key.id

  secret_data = base64encode(random_string.jwt_signing_key.result)
}

### External API services ###

# Sendgrid token

variable "SENDGRID_API_TOKEN" {}

resource "google_secret_manager_secret" "sendgrid_api_token" {
  secret_id = "sendgrid_api_token"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "sendgrid_api_token" {
  secret = google_secret_manager_secret.sendgrid_api_token.id

  secret_data = var.SENDGRID_API_TOKEN
}

# Segment API token

variable "SEGMENT_API_TOKEN" {}

resource "google_secret_manager_secret" "segment_api_token" {
  secret_id = "segment_api_token"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "segment_api_token" {
  secret = google_secret_manager_secret.segment_api_token.id

  secret_data = var.SEGMENT_API_TOKEN
}

# PostHog API token

variable "POSTHOG_API_KEY" {}

resource "google_secret_manager_secret" "posthog_api_key" {
  secret_id = "posthog_api_key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "posthog_api_key" {
  secret = google_secret_manager_secret.posthog_api_key.id

  secret_data = var.POSTHOG_API_KEY
}

# PostHog personal API key

variable "POSTHOG_PERSONAL_API_KEY" {}

resource "google_secret_manager_secret" "posthog_personal_api_key" {
  secret_id = "posthog_personal_api_key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "posthog_personal_api_key" {
  secret = google_secret_manager_secret.posthog_personal_api_key.id

  secret_data = var.POSTHOG_PERSONAL_API_KEY
}

# Algolia app ID

resource "google_secret_manager_secret" "algolia_application_id" {
  secret_id = "algolia_application_id"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "algolia_application_id" {
  secret = google_secret_manager_secret.algolia_application_id.id

  secret_data = var.ALGOLIA_APPLICATION_ID
}

# Algolia API key

resource "google_secret_manager_secret" "algolia_api_key" {
  secret_id = "algolia_api_key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "algolia_api_key" {
  secret = google_secret_manager_secret.algolia_api_key.id

  secret_data = var.ALGOLIA_API_KEY
}

# Google SSO Client ID

resource "google_secret_manager_secret" "google_sso_client_id" {
  secret_id = "google_sso_client_id"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "google_sso_client_id" {
  secret = google_secret_manager_secret.google_sso_client_id.id

  secret_data = var.GOOGLE_SSO_OAUTH2_CLIENT_ID
}

# Google SSO Client Secret

resource "google_secret_manager_secret" "google_sso_client_secret" {
  secret_id = "google_sso_client_secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "google_sso_client_secret" {
  secret = google_secret_manager_secret.google_sso_client_secret.id

  secret_data = var.GOOGLE_SSO_OAUTH2_CLIENT_SECRET
}
