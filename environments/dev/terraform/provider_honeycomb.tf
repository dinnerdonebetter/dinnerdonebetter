variable "HONEYCOMB_API_KEY" {}

# Configure the AWS Provider
provider "honeycombio" {
  api_key = var.HONEYCOMB_API_KEY
}

data "aws_ssm_parameter" "honeycomb_api_key" {
  name = "PRIXFIXE_HONEYCOMB_API_KEY"
}