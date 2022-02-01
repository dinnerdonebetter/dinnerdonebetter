resource "digitalocean_database_cluster" "database" {
  name       = "prixfixe-dev"
  engine     = "pg"
  version    = "13"
  size       = "db-s-1vcpu-1gb"
  region     = local.region
  node_count = 1
}