data "digitalocean_kubernetes_versions" "example" {
  version_prefix = "1.18."
}

resource "digitalocean_kubernetes_cluster" "dev" {
  name         = "dev"
  region       = local.region
  auto_upgrade = true
  version      = data.digitalocean_kubernetes_versions.example.latest_version

  maintenance_policy {
    start_time  = "04:00"
    day         = "sunday"
  }

  node_pool {
    name       = "default"
    size       = "s-1vcpu-1gb"
    node_count = 3
  }
}