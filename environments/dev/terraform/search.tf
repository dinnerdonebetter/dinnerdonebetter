resource "aws_elasticsearch_domain" "search" {
  domain_name           = "dev-search"
  elasticsearch_version = "2.3"

  cluster_config {
    instance_type = "t2.micro.elasticsearch"
  }

  tags = merge(var.default_tags, {})
}

resource "aws_ssm_parameter" "search_url" {
  name  = "PRIXFIXE_ELASTICSEARCH_INSTANCE_URL"
  type  = "String"
  value = aws_elasticsearch_domain.search.arn
}