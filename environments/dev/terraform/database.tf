resource "digitalocean_database_cluster" "database" {
  name                 = "database"
  engine               = "pg"
  version              = "13"
  size                 = "db-s-1vcpu-1gb"
  region               = local.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.dev.id
}

resource "digitalocean_database_firewall" "firewall" {
  cluster_id = digitalocean_database_cluster.database.id

  rule {
    type  = "k8s"
    value = digitalocean_kubernetes_cluster.dev.id
  }
}

resource "digitalocean_database_db" "api_database" {
  cluster_id = digitalocean_database_cluster.database.id
  name       = "prixfixe"
}

resource "digitalocean_database_user" "api_user" {
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

resource "kubernetes_secret_v1" "database_credentials" {
  metadata {
    namespace = local.kubernetes_namespace
    name      = "config.database"
  }

  data = {
    "connection_string" = format(
      "user=%s dbname=%s password='%s' host=%s port=%s",
      digitalocean_database_user.api_user.name,
      digitalocean_database_db.api_database.name,
      digitalocean_database_user.api_user.password,
      digitalocean_database_cluster.database.host,
      digitalocean_database_cluster.database.port,
    )
  }
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