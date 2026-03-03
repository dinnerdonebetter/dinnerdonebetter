locals {
  k8s_namespace = "prod"
}

provider "kubernetes" {
  config_path    = "./terraform_kubeconfig"
  config_context = "${local.k8s_namespace}_context"
}

data "google_container_cluster" "prod_cluster" {
  name     = "prod"
  location = local.gcp_region
}

resource "kubernetes_namespace" "prod" {
  metadata {
    name = local.k8s_namespace
    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [data.google_container_cluster.prod_cluster]
}

# Kubernetes secrets

# APNs .p8 key mounted as file for push notifications (async message handler)
resource "kubernetes_secret" "apns_credentials" {
  metadata {
    name      = "apns-credentials"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    # Key becomes filename when mounted; value is p8 content (Terraform base64-encodes automatically)
    "apns-auth-key.p8" = var.APNS_AUTH_KEY_P8
  }
}

resource "kubernetes_secret" "cloudflare_api_key" {
  metadata {
    name      = "cloudflare-api-key"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    "token" = var.CLOUDFLARE_API_TOKEN
  }
}

resource "kubernetes_config_map_v1" "pubsub_topics" {
  metadata {
    name      = "pubsub-topic-names"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    data_changes               = google_pubsub_topic.data_changes_topic.id
    outbound_emails            = google_pubsub_topic.outbound_emails_topic.id
    search_index_requests      = google_pubsub_topic.search_index_requests_topic.id
    mobile_notifications       = google_pubsub_topic.mobile_notifications_topic.id
    user_data_aggregator       = google_pubsub_topic.user_data_aggregator_topic.id
    webhook_execution_requests = google_pubsub_topic.webhook_execution_requests_topic.id
  }
}

resource "kubernetes_secret" "api_service_config" {
  metadata {
    name      = "api-service-config"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    OAUTH2_TOKEN_ENCRYPTION_KEY     = random_string.oauth2_token_encryption_key.result
    JWT_SIGNING_KEY                 = base64encode(random_string.jwt_signing_key.result)
    DATABASE_HOST                   = google_sql_database_instance.prod.private_ip_address
    DATABASE_USERNAME               = local.api_database_username
    DATABASE_PASSWORD               = random_password.api_user_database_password.result
    SENDGRID_API_TOKEN              = var.SENDGRID_API_KEY
    SEGMENT_API_TOKEN               = var.SEGMENT_API_TOKEN
    POSTHOG_API_KEY                 = var.POSTHOG_API_KEY
    POSTHOG_PERSONAL_API_KEY        = var.POSTHOG_PERSONAL_API_KEY
    ALGOLIA_APPLICATION_ID          = var.ALGOLIA_APPLICATION_ID
    ALGOLIA_API_KEY                 = var.ALGOLIA_API_KEY
    GOOGLE_SSO_OAUTH2_CLIENT_ID     = var.GOOGLE_SSO_OAUTH2_CLIENT_ID
    GOOGLE_SSO_OAUTH2_CLIENT_SECRET = var.GOOGLE_SSO_OAUTH2_CLIENT_SECRET
    PUSH_NOTIFICATIONS_APNS_KEY_ID  = var.APNS_KEY_ID
  }
}

# this is the sort of resource that should probably ideally live in the infra folder, but it's here for now
# because I haven't yet wanted to fuss with figuring out how to manage the code that creates the cluster
# alongside the code that creates resources in that cluster.
# Admin webapp and MCP server: OAuth2 credentials + cookie config
# Maps to env vars DINNER_DONE_BETTER_API_SERVICE_OAUTH2_API_CLIENT_ID / _SECRET
resource "kubernetes_secret" "admin_webapp_config" {
  metadata {
    name      = "dinner-done-better-admin-webapp-config"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    OAUTH2_CLIENT_ID     = var.ADMIN_WEBAPP_OAUTH2_CLIENT_ID
    OAUTH2_CLIENT_SECRET = var.ADMIN_WEBAPP_OAUTH2_CLIENT_SECRET
    COOKIE_NAME          = var.ADMIN_WEBAPP_COOKIE_NAME
    COOKIE_HASH_KEY      = base64encode(random_string.admin_webapp_cookie_hash_key.result)
    COOKIE_BLOCK_KEY     = base64encode(random_string.admin_webapp_cookie_block_key.result)
    COOKIE_DOMAIN        = var.ADMIN_WEBAPP_COOKIE_DOMAIN
  }
}

# MCP server OAuth2 credentials for DINNER_DONE_BETTER_API_SERVICE_OAUTH2_API_CLIENT_ID / _SECRET
resource "kubernetes_secret" "mcp_server_config" {
  metadata {
    name      = "mcp-server-config"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    OAUTH2_CLIENT_ID     = var.MCP_SERVICE_OAUTH2_CLIENT_ID
    OAUTH2_CLIENT_SECRET = var.MCP_SERVICE_OAUTH2_CLIENT_SECRET
  }
}

# this is the sort of resource that should probably ideally live in the infra folder, but it's here for now
# because I haven't yet wanted to fuss with figuring out how to manage the code that creates the cluster
# alongside the code that creates resources in that cluster.
resource "kubernetes_secret" "grafana_cloud_creds" {
  metadata {
    name      = "grafana-cloud-creds"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [kubernetes_namespace.prod]

  data = {
    GRAFANA_CLOUD_PROMETHEUS_USERNAME = var.GRAFANA_CLOUD_PROMETHEUS_USERNAME
    GRAFANA_CLOUD_PROMETHEUS_PASSWORD = var.GRAFANA_CLOUD_PROMETHEUS_PASSWORD
    GRAFANA_CLOUD_LOKI_USERNAME       = var.GRAFANA_CLOUD_LOKI_USERNAME
    GRAFANA_CLOUD_LOKI_PASSWORD       = var.GRAFANA_CLOUD_LOKI_PASSWORD
    GRAFANA_CLOUD_TEMPO_USERNAME      = var.GRAFANA_CLOUD_TEMPO_USERNAME
    GRAFANA_CLOUD_TEMPO_PASSWORD      = var.GRAFANA_CLOUD_TEMPO_PASSWORD
  }
}
