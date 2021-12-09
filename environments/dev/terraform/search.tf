data "ec_stack" "latest" {
  version_regex = "latest"
  region        = "us-east-1"
}

resource "ec_deployment" "dev" {
  name = "dev_api"

  # Mandatory fields
  region                 = "us-east-1"
  version                = data.ec_stack.latest.version
  deployment_template_id = "aws-storage-optimized"

  elasticsearch {
    autoscale = false
  }
}

resource "aws_ssm_parameter" "search_url" {
  name  = "PRIXFIXE_ELASTICSEARCH_INSTANCE_URL"
  type  = "String"
  value = ec_deployment.dev.elasticsearch[0].https_endpoint
}

resource "aws_ssm_parameter" "search_url" {
  name  = "PRIXFIXE_ELASTICSEARCH_USERNAME"
  type  = "String"
  value = ec_deployment.dev.elasticsearch_username
}

resource "aws_ssm_parameter" "search_url" {
  name  = "PRIXFIXE_ELASTICSEARCH_PASSWORD"
  type  = "String"
  value = ec_deployment.dev.elasticsearch_password
}
