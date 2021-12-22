locals {
  redis_username = "dev-api-server"
}

resource "aws_elasticache_parameter_group" "dev_api" {
  name   = "dev-api-params"
  family = "redis6.x"
}

resource "aws_security_group" "redis" {
  name        = "redis"
  description = "Redis access"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port        = 6379
    to_port          = 6379
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_elasticache_subnet_group" "redis" {
  name       = "redis-subnet"
  subnet_ids = [for x in aws_subnet.private_subnets : x.id]
}

resource "aws_elasticache_cluster" "dev_api" {
  cluster_id           = "dev-api"
  engine               = "redis"
  node_type            = "cache.t4g.micro"
  num_cache_nodes      = 1
  parameter_group_name = aws_elasticache_parameter_group.dev_api.name
  engine_version       = "6.x"
  port                 = 6379

  subnet_group_name = aws_elasticache_subnet_group.redis.name
  security_group_ids = [
    aws_security_group.redis.id,
  ]
}

resource "random_password" "redis_password" {
  length           = 64
  special          = true
  override_special = "-._~()'!*"
}

resource "aws_elasticache_user" "dev_api" {
  user_id       = "dev-api"
  user_name     = local.redis_username
  access_string = "on ~* +@pubsub"
  engine        = "REDIS"
  passwords     = [random_password.redis_password.result]
}

resource "aws_ssm_parameter" "pubsub_server_username" {
  name  = "PRIXFIXE_PUBSUB_SERVER_USERNAME"
  type  = "String"
  value = local.redis_username
}

resource "aws_ssm_parameter" "pubsub_server_password" {
  name  = "PRIXFIXE_PUBSUB_SERVER_PASSWORD"
  type  = "String"
  value = random_password.redis_password.result
}

resource "aws_ssm_parameter" "pubsub_server_url" {
  name  = "PRIXFIXE_PUBSUB_SERVER_URLS"
  type  = "String"
  value = join(",", [for x in aws_elasticache_cluster.dev_api.cache_nodes : format("%s:6379", x.address)])
}
