variable "CLOUDFLARE_API_KEY" {}
variable "CLOUDFLARE_ZONE_ID" {}


provider "cloudflare" {
  api_token = var.CLOUDFLARE_API_KEY
}