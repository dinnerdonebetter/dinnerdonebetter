locals {
  k8s_namespace = "dev"
}

provider "kubernetes" {
  config_path    = "./terraform_kubeconfig"
  config_context = "${local.k8s_namespace}_context"
}

data "google_container_cluster" "dev_cluster" {
  name     = "dev"
  location = local.gcp_region
}

# Kubernetes secrets

variable "API_SERVICE_OAUTH2_CLIENT_ID" {}
variable "API_SERVICE_OAUTH2_CLIENT_SECRET" {}

resource "random_bytes" "webapp_cookie_encryption_key" {
  length = 32
}

resource "random_bytes" "cookie_encryption_iv" {
  length = 32
}

resource "kubernetes_secret" "cloudflare_api_key" {
  metadata {
    name      = "webapp-config"
    namespace = local.k8s_namespace

    annotations = {
      "managed_by" = "terraform"
    }

    labels = {
      "managed_by" = "terraform"
    }
  }

  depends_on = [data.google_container_cluster.dev_cluster]

  data = {
    COOKIE_ENCRYPTION_KEY                   = random_bytes.webapp_cookie_encryption_key.hex
    BASE64_COOKIE_ENCRYPT_IV                = random_bytes.cookie_encryption_iv.base64
    API_SERVICE_OAUTH2_CLIENT_ID            = var.API_SERVICE_OAUTH2_CLIENT_ID
    API_SERVICE_OAUTH2_CLIENT_SECRET        = var.API_SERVICE_OAUTH2_CLIENT_SECRET
    SEGMENT_API_TOKEN                       = var.SEGMENT_API_TOKEN
    DINNER_DONE_BETTER_OAUTH2_CLIENT_ID     = var.DINNER_DONE_BETTER_OAUTH2_CLIENT_ID
    DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET = var.DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET
  }
}
