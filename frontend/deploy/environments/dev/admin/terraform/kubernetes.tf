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

resource "kubernetes_config_map_v1" "frontend_admin_app_service_configmap" {
  metadata {
    name      = "frontend-admin-app-service-config"
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
    APIEndpoint : "dinner-done-better.dev.svc.cluster.local:8000",
    CookieEncryptionKey : random_string.cookie_encryption_key.result,
    CookieEncryptionIV : random_bytes.cookie_encryption_iv.base64,
    APIOAuth2ClientID : var.DINNER_DONE_BETTER_OAUTH2_CLIENT_ID,
    APIOAuth2ClientSecret : var.DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET
  }
}
