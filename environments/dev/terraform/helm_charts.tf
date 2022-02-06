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
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "external-dns"

  set {
    name  = "logFormat"
    value = "json"
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

#resource "helm_release" "nats" {
#  name       = "nats"
#  repository = "https://nats-io.github.io/k8s/helm/charts"
#  chart      = "nats"
#
#  #  set {
#  #    name = "jetstream.enabled"
#  #    value = true
#  #  }
#}

#resource "helm_release" "prometheus" {
#  name       = "prometheus"
#  repository = "https://prometheus-community.github.io/helm-charts"
#  chart      = "prometheus"
#
#  #  set {
#  #    name = "ingress.annotations"
#  #    value = jsonencode({
#  #      things : "stuff"
#  #    })
#  #  }
#}

#resource "helm_release" "grafana" {
#  name       = "grafana"
#  repository = "https://grafana.github.io/helm-charts"
#  chart      = "grafana"
#
#  #  set {
#  #    name = "ingress.annotations"
#  #    value = jsonencode({
#  #      things : "stuff"
#  #    })
#  #  }
#}
