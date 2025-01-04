locals {
  public_url = "dinnerdonebetter.dev"
}

resource "sendgrid_domain_authentication" "dev" {
  domain             = local.public_url
  subdomain          = "email"
  is_default         = true
  automatic_security = false
}

resource "sendgrid_link_branding" "default" {
  domain     = local.public_url
  is_default = true
}

resource "cloudflare_record" "domain_verification_records" {
  for_each = {
    "0" : sendgrid_domain_authentication.dev.dns[0],
    "1" : sendgrid_domain_authentication.dev.dns[1],
    "2" : sendgrid_domain_authentication.dev.dns[2],
  }

  zone_id  = var.CLOUDFLARE_ZONE_ID
  name     = each.value.host
  type     = upper(each.value.type)
  value    = each.value.data
  ttl      = 1
  proxied  = false
  comment  = "Managed by Terraform"
  priority = upper(each.value.type) == "MX" ? 10 : null
}
