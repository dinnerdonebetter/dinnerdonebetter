locals {
  kubernetes_namespace = "dev"
}

resource "digitalocean_kubernetes_cluster" "dev" {
  name         = "dev"
  region       = local.region
  auto_upgrade = true
  version      = "1.21.9-do.0"
  vpc_uuid     = digitalocean_vpc.dev.id

  maintenance_policy {
    start_time = "04:00"
    day        = "sunday"
  }

  node_pool {
    name       = "default"
    size       = "s-2vcpu-4gb"
    node_count = 2
  }
}

resource "digitalocean_project_resources" "dev_cluster" {
  project = digitalocean_project.prixfixe_dev.id
  resources = [
    # https://www.digitalocean.com/community/questions/attach-kubernetes-cluster-to-project
    format("do:kubernetes:%s", digitalocean_kubernetes_cluster.dev.id),
  ]
}

provider "kubernetes" {
  host  = digitalocean_kubernetes_cluster.dev.endpoint
  token = digitalocean_kubernetes_cluster.dev.kube_config[0].token
  cluster_ca_certificate = base64decode(
    digitalocean_kubernetes_cluster.dev.kube_config[0].cluster_ca_certificate
  )
}

resource "kubernetes_namespace_v1" "dev_namespace" {
  metadata {
    name = local.kubernetes_namespace
  }
}

provider "helm" {
  kubernetes {
    host  = digitalocean_kubernetes_cluster.dev.endpoint
    token = digitalocean_kubernetes_cluster.dev.kube_config[0].token
    cluster_ca_certificate = base64decode(
      digitalocean_kubernetes_cluster.dev.kube_config[0].cluster_ca_certificate
    )
  }
}

resource "helm_release" "kubewatch" {
  name       = "kubewatch"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "external-dns"

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