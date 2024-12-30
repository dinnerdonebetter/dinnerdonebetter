locals {
  k8s_namespace = "dev"
}

variable "KUBECONFIG" {
  type = string
}

output "kubeconfig" {
  value = var.KUBECONFIG
}

# provider "kubernetes" {
#   config_path    = var.KUBECONFIG
#   config_context = "${local.k8s_namespace}_context"
# }

# resource "kubernetes_namespace" "dev" {
#   metadata {
#     annotations = {
#       exampleAnnotation = "example-annotation"
#     }
#
#     labels = {
#       exampleLabel = "example-label"
#     }
#
#     name = local.k8s_namespace
#   }
# }

# # Kubernetes secrets
#
# resource "kubernetes_secret" "cloudflare_api_key" {
#   metadata {
#     name      = "cloudflare-api-key"
#     namespace = local.k8s_namespace
#   }
#
#   data = {
#     "token" = var.CLOUDFLARE_API_TOKEN
#   }
# }
#
# resource "kubernetes_secret" "pubsub_topics" {
#   metadata {
#     name      = "pubsub-topic-names"
#     namespace = local.k8s_namespace
#   }
#
#   data = {
#     "data_changes"               = google_pubsub_topic.data_changes_topic.name
#     "outbound_emails"            = google_pubsub_topic.outbound_emails_topic.name
#     "search_index_requests"      = google_pubsub_topic.search_index_requests_topic.name
#     "user_data_aggregator"       = google_pubsub_topic.user_data_aggregator_topic.name
#     "webhook_execution_requests" = google_pubsub_topic.webhook_execution_requests_topic.name
#   }
# }
#
# resource "kubernetes_secret" "api_service_config" {
#   metadata {
#     name      = "api-service-config"
#     namespace = local.k8s_namespace
#   }
#
#   data = {
#     "api-service-config.json"         = "${file("${path.module}/service-config.json")}"
#     "OAUTH2_TOKEN_ENCRYPTION_KEY"     = random_string.oauth2_token_encryption_key.result
#     "JWT_SIGNING_KEY"                 = base64encode(random_string.jwt_signing_key.result)
#     "SENDGRID_API_TOKEN"              = var.SENDGRID_API_TOKEN
#     "SEGMENT_API_TOKEN"               = var.SEGMENT_API_TOKEN
#     "POSTHOG_API_KEY"                 = var.POSTHOG_API_KEY
#     "POSTHOG_PERSONAL_API_KEY"        = var.POSTHOG_PERSONAL_API_KEY
#     "ALGOLIA_APPLICATION_ID"          = var.ALGOLIA_APPLICATION_ID
#     "ALGOLIA_API_KEY"                 = var.ALGOLIA_API_KEY
#     "GOOGLE_SSO_OAUTH2_CLIENT_ID"     = var.GOOGLE_SSO_OAUTH2_CLIENT_ID
#     "GOOGLE_SSO_OAUTH2_CLIENT_SECRET" = var.GOOGLE_SSO_OAUTH2_CLIENT_SECRET
#   }
# }
