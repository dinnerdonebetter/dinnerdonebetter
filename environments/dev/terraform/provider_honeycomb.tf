variable "HONEYCOMB_API_KEY" {}

# Configure the AWS Provider
provider "honeycombio" {
  api_key = var.HONEYCOMB_API_KEY
}

data "aws_ssm_parameter" "honeycomb_api_key" {
  name = "PRIXFIXE_HONEYCOMB_API_KEY"
}

resource "honeycombio_marker" "server_logs_deploy_marker" {
  type = "deploy"

  dataset = honeycombio_dataset.dev_api_server_logs.name
}

resource "honeycombio_marker" "worker_logs_deploy_marker" {
  type = "deploy"

  dataset = honeycombio_dataset.dev_worker_logs.name
}

resource "honeycombio_marker" "postgres_logs_deploy_marker" {
  type = "deploy"

  dataset = honeycombio_dataset.dev_postgres_rds_logs.name
}