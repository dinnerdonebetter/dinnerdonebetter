provider "helm" {
  kubernetes {
    host  = digitalocean_kubernetes_cluster.dev.endpoint
    token = digitalocean_kubernetes_cluster.dev.kube_config[0].token
    cluster_ca_certificate = base64decode(
      digitalocean_kubernetes_cluster.dev.kube_config[0].cluster_ca_certificate
    )
  }
}

resource "helm_release" "external_dns" {
  name       = "external-dns"
  repository = "https://kubernetes-sigs.github.io/external-dns"
  chart      = "external-dns"

  set {
    name  = "logFormat"
    value = "json"
  }

  set {
    name  = "namespace"
    value = local.kubernetes_namespace
  }

  set_sensitive {
    name  = "cloudflare.apiToken"
    value = var.CLOUDFLARE_API_TOKEN
  }

  set {
    name  = "provider"
    value = "cloudflare"
  }
}

#resource "helm_release" "prometheus" {
#  name       = "prometheus"
#  repository = "https://prometheus-community.github.io/helm-charts"
#  chart      = "prometheus-community"
#
#  set {
#    name  = "ingress.annotations"
#    value = {
#      things = "stuff"
#    }
#  }
#
#  set {
#    name = "namespaceOverride"
#    value = local.kubernetes_namespace
#  }
#}
#
#resource "helm_release" "grafana" {
#  name       = "grafana"
#  repository = "https://grafana.github.io/helm-charts"
#  chart      = "grafana"
#
#  set {
#    name  = "ingress.annotations"
#    value = {
#      things = "stuff"
#    }
#  }
#
#  set {
#    name = "namespaceOverride"
#    value = local.kubernetes_namespace
#  }
#}
