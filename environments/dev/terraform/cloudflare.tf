variable "cf_token" {}
variable "cf_zone_id" {}

provider "cloudflare" {
  version   = "~> 2.0"
  api_token = var.cf_token
}

resource "cloudflare_record" "db_record" {
  zone_id = var.cf_zone_id
  name    = "database"
  value   = data.digitalocean_database_cluster.database.host
  type    = "CNAME"
  ttl     = 3600
}

resource "cloudflare_record" "root_record_v4" {
  zone_id = var.cf_zone_id
  name    = "@"
  value   = data.digitalocean_droplet.dev_server.ipv4_address
  type    = "A"
  ttl     = 3600
}

resource "cloudflare_record" "root_record_v6" {
  zone_id = var.cf_zone_id
  name    = "@"
  value   = data.digitalocean_droplet.dev_server.ipv6_address
  type    = "AAAA"
  ttl     = 3600
}
