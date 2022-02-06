resource "kubernetes_deployment" "api_server" {
  metadata {
    name      = "prixfixe-api-server"
    namespace = local.kubernetes_namespace
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "api-server"
      }
    }

    template {
      metadata {
        name = "prixfixe-api-server"
        labels = {
          app = "api-server"
        }
      }

      spec {
        image_pull_secrets {
          name = "registry-dev-prixfixe"
        }
        container {
          name              = "api-server-container"
          image             = "registry.digitalocean.com/dev-prixfixe/api_server:latest"
          image_pull_policy = "Always"

          env {
            name  = "CONFIGURATION_FILEPATH"
            value = "/config/service-config.json"
          }

          env {
            name = "DATABASE_CONFIGURATION"
            value_from {
              secret_key_ref {
                name = "config.database"
                key  = "connection_string"
              }
            }
          }

          env {
            name = "PRIXFIXE_COOKIE_HASH_KEY"
            value_from {
              secret_key_ref {
                name = "config.auth"
                key  = "cookie_hash_key"
              }
            }
          }

          env {
            name = "PRIXFIXE_COOKIE_BLOCK_KEY"
            value_from {
              secret_key_ref {
                name = "config.auth"
                key  = "cookie_block_key"
              }
            }
          }

          env {
            name = "PRIXFIXE_PASETO_LOCAL_KEY"
            value_from {
              secret_key_ref {
                name = "config.auth"
                key  = "paseto_local_key"
              }
            }
          }

          env {
            name = "PRIXFIXE_SENDGRID_API_TOKEN"
            value_from {
              secret_key_ref {
                name = "config.third-party.sendgrid"
                key  = "sendgrid_api_token"
              }
            }
          }

          env {
            name = "PRIXFIXE_SEGMENT_API_TOKEN"
            value_from {
              secret_key_ref {
                name = "config.third-party.segment"
                key  = "segment_api_token"
              }
            }
          }

          resources {
            requests = {
              cpu    = "64Mi"
              memory = "250m"
            }
            limits = {
              cpu    = "256Mi"
              memory = "500m"
            }
          }

          port {
            container_port = 8000
          }

          volume_mount {
            name       = "config"
            mount_path = "/config/"
            read_only  = true
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_namespace_v1.dev_namespace,
    kubernetes_secret_v1.config_file,
    kubernetes_secret_v1.config_auth,
    kubernetes_secret_v1.config_sendgrid,
    kubernetes_secret_v1.config_segment,
  ]
}

resource "kubernetes_service_v1" "api_service" {
  metadata {
    name      = "prixfixe-api-service"
    namespace = local.kubernetes_namespace
    labels = {
      app = "api-server"
    }
    annotations = {
      "service.beta.kubernetes.io/do-loadbalancer-id" : digitalocean_loadbalancer.public.id
      "service.beta.kubernetes.io/do-loadbalancer-name" : "api.prixfixe.dev"
      "service.beta.kubernetes.io/do-loadbalancer-protocol" : "http"
      "service.beta.kubernetes.io/do-loadbalancer-hostname" : "api.prixfixe.dev"
      "service.beta.kubernetes.io/do-loadbalancer-size-unit" : 1
      "external-dns.alpha.kubernetes.io/hostname" : "api.prixfixe.dev"
      "external-dns.alpha.kubernetes.io/ttl" : 120
      "external-dns.alpha.kubernetes.io/cloudflare-proxied" : true
    }
  }
  spec {
    type = "LoadBalancer"
    selector = {
      app = "api-server"
    }

    port {
      name        = "http"
      protocol    = "TCP"
      port        = 8000
      target_port = 80
    }
  }

  depends_on = [
    digitalocean_loadbalancer.public,
    kubernetes_namespace_v1.dev_namespace,
  ]
}