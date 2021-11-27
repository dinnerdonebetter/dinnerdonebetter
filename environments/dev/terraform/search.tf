resource "aws_elasticsearch_domain" "search" {
  domain_name           = "dev-search"
  elasticsearch_version = "7.10"

  cluster_config {
    instance_type = "t2.micro.search"
  }

  tags = merge(var.default_tags, {})
}

resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_ELASTICSEARCH_INSTANCE_URL"
  type  = "String"
  value = aws_elasticsearch_domain.search.arn
}