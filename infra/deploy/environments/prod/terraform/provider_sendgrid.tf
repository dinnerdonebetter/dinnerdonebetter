variable "SENDGRID_API_KEY" {}

provider "sendgrid" {
  api_key = var.SENDGRID_API_KEY
}
