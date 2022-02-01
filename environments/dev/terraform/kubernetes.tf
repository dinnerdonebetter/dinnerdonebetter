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

output "kubernetes_host" {
  value = digitalocean_kubernetes_cluster.dev.endpoint
}

output "kubernetes_token" {
  value = digitalocean_kubernetes_cluster.dev.kube_config[0].token
}

output "kubernetes_cluster_cert" {
  value = base64decode(
    digitalocean_kubernetes_cluster.dev.kube_config[0].cluster_ca_certificate
  )
}