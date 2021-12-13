variable "HONEYCOMB_API_KEY" {}

# Configure the AWS Provider
provider "honeycombio" {
  api_key = var.HONEYCOMB_API_KEY
}