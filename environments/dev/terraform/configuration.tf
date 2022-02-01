resource "digitalocean_spaces_bucket" "config" {
  name          = "dev-config"
  region        = local.region
  acl           = "private"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "configuration_file" {
  region       = digitalocean_spaces_bucket.config.region
  bucket       = digitalocean_spaces_bucket.config.name
  key          = "dev.api.json"
  content      = file("${path.module}/service-config.json")
  content_type = "application/json"
}