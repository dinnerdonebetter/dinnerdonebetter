locals {
  database_username = "prixfixe_api"
}

resource "random_password" "database_password" {
  length           = 64
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = [for x in aws_subnet.public_subnets : x.id]
}

resource "aws_rds_cluster" "api_database" {
  cluster_identifier = "api-database"
  database_name      = "prixfixe"
  engine             = "aurora-postgresql"

  engine_mode = "serverless"
  scaling_configuration {
    auto_pause               = true
    min_capacity             = 2
    max_capacity             = 2
    seconds_until_auto_pause = 300
    timeout_action           = "ForceApplyCapacityChange"
  }

  # master_username              = local.database_username
  # master_password              = random_password.database_password.result
  master_username              = "prixfixe_api"
  master_password              = "password"
  preferred_backup_window      = "05:00-08:00"
  preferred_maintenance_window = "sat:01:00-sat:04:00"
  apply_immediately            = true
  enable_http_endpoint         = true
  storage_encrypted            = true
  skip_final_snapshot          = true
  copy_tags_to_snapshot        = true
  backup_retention_period      = 7

  db_subnet_group_name = aws_db_subnet_group.default.name
  vpc_security_group_ids = [
    aws_security_group.allow_postgres.id,
  ]
}

resource "aws_ssm_parameter" "database_url" {
  name  = "PRIXFIXE_DATABASE_URL"
  type  = "String"
  value = format("postgres://%s:%s@%s:%d/prixfixe", local.database_username, random_password.database_password.result, aws_rds_cluster.api_database.endpoint, aws_rds_cluster.api_database.port)
}
