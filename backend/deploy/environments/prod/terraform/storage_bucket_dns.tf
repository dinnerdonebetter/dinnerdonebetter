# CNAME records for domain-named GCS buckets (HTTP website hosting).
# Points media.dinnerdonebetter.com and userdata.dinnerdonebetter.com to GCS.
# proxied = true: Cloudflare terminates TLS and proxies to GCS over HTTP.
resource "cloudflare_dns_record" "media_storage_bucket" {
  zone_id = var.CLOUDFLARE_ZONE_ID
  name    = "media"
  type    = "CNAME"
  content = "c.storage.googleapis.com"
  proxied = true
  ttl     = 1
  comment = "GCS bucket media.dinnerdonebetter.com (HTTPS via Cloudflare)"
}

resource "cloudflare_dns_record" "userdata_storage_bucket" {
  zone_id = var.CLOUDFLARE_ZONE_ID
  name    = "userdata"
  type    = "CNAME"
  content = "c.storage.googleapis.com"
  proxied = true
  ttl     = 1
  comment = "GCS bucket userdata.dinnerdonebetter.com (HTTPS via Cloudflare)"
}
