locals {
  redis_username = "dev-api-server"
}

resource "aws_elasticache_parameter_group" "dev_api" {
  name   = "dev-api-params"
  family = "redis6.x"
}

resource "aws_elasticache_cluster" "dev_api" {
  cluster_id           = "dev-api"
  engine               = "redis"
  node_type            = "cache.t4g.micro"
  num_cache_nodes      = 1
  parameter_group_name = aws_elasticache_parameter_group.dev_api.name
  engine_version       = "6.x"
  port                 = 6379
}

resource "random_password" "redis_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "aws_elasticache_user" "dev_api" {
  user_id       = "dev-api"
  user_name     = local.redis_username
  access_string = "on ~* +@pubsub"
  engine        = "REDIS"
  passwords     = [random_password.database_password.result]
}

output "es_cluster" {
  value = aws_elasticache_cluster.dev_api
}

resource "aws_ssm_parameter" "pubsub_server_url" {
  name  = "PRIXFIXE_PUBSUB_SERVER_URL"
  type  = "String"
  value = aws_elasticache_cluster.dev_api.arn
}

resource "aws_ssm_parameter" "pubsub_server_username" {
  name  = "PRIXFIXE_PUBSUB_SERVER_USERNAME"
  type  = "String"
  value = local.redis_username
}

resource "aws_ssm_parameter" "pubsub_server_password" {
  name  = "PRIXFIXE_PUBSUB_SERVER_PASSWORD"
  type  = "String"
  value = random_password.database_password.result
}
