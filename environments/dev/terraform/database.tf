resource "digitalocean_database_cluster" "database" {
  name                 = "database"
  engine               = "pg"
  version              = "13"
  size                 = "db-s-1vcpu-1gb"
  region               = local.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.dev.id
}

resource "digitalocean_database_user" "user-example" {
  cluster_id = digitalocean_database_cluster.database.id
  name       = "prixfixe_api"
}


resource "digitalocean_project_resources" "dev_db" {
  project = digitalocean_project.prixfixe_dev.id
  resources = [
    # https://github.com/digitalocean/api-v2/issues/179
    format("do:dbaas:%s", digitalocean_database_cluster.database.id),
  ]
}

resource "cloudflare_record" "database_dot_prixfixe_dot_dev" {
  zone_id         = var.CLOUDFLARE_ZONE_ID
  name            = "database.prixfixe.dev"
  value           = digitalocean_database_cluster.database.host
  type            = "CNAME"
  proxied         = true
  allow_overwrite = true
  ttl             = 1
}
