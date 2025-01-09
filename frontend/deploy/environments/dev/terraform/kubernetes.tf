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

resource "kubernetes_secret" "frontend_service_config" {
  metadata {
    name      = "frontend-service-config"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [data.google_container_cluster.dev_cluster]

  data = {
    "token" = var.CLOUDFLARE_API_TOKEN
  }
}

resource "kubernetes_secret" "api_service_config" {
  metadata {
    name      = "frontend-service-config"
    namespace = local.k8s_namespace

    annotations = {
      (local.managed_by_label) = "terraform"
    }

    labels = {
      (local.managed_by_label) = "terraform"
    }
  }

  depends_on = [data.google_container_cluster.dev_cluster]

  data = {
    OAUTH2_CLIENT_ID     = var.DINNER_DONE_BETTER_OAUTH2_CLIENT_ID
    OAUTH2_CLIENT_SECRET = var.DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET
    COOKIE_ENCRYPTION_KEY = base64encode(random_string.cookie_encryption_key.result)
    COOKIE_ENCRYPTION_IV = random_bytes.cookie_encryption_iv.base64
    SEGMENT_API_TOKEN    = var.SEGMENT_API_TOKEN
  }
}
