resource "digitalocean_kubernetes_cluster" "dev" {
  name         = "dev"
  region       = local.region
  auto_upgrade = true
  version      = "1.21.9-do.0"

  maintenance_policy {
    start_time = "04:00"
    day        = "sunday"
  }

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 3
  }
}

resource "digitalocean_project_resources" "dev_cluster" {
  project = digitalocean_project.prixfixe_dev.id
  resources = [
    format("do:kubernetes:%s", digitalocean_kubernetes_cluster.dev.id),
  ]
}

output "k8s_cluster_endpoint" {
  value = digitalocean_kubernetes_cluster.dev.endpoint
}

provider "kubernetes" {
  host  = digitalocean_kubernetes_cluster.dev.endpoint
  token = digitalocean_kubernetes_cluster.dev.kube_config[0].token
  cluster_ca_certificate = base64decode(
    digitalocean_kubernetes_cluster.dev.kube_config[0].cluster_ca_certificate
  )
}