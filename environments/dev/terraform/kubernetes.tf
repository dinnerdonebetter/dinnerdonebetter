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

output "k8s_cluster_enpdoint" {
  value = digitalocean_kubernetes_cluster.dev.endpoint
}
