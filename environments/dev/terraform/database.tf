locals {
  database_username = "api"
}

resource "random_password" "database_password" {
  length  = 64
  special = true
}

resource "aws_rds_cluster" "api_database" {
  cluster_identifier      = "dev_database"
  engine                  = "aurora-postgresql"
  availability_zones      = ["us-east-1"]
  database_name           = "prixfixe"
  engine_mode             = "serverless"
  master_username         = local.database_username
  master_password         = random_password.database_password.result
  backup_retention_period = 7
  storage_encrypted       = true
  preferred_backup_window = "01:00-05:00"
}

resource "aws_ssm_parameter" "database_username" {
  name  = "database_password"
  type  = "String"
  value = local.database_username
}

resource "aws_ssm_parameter" "database_password" {
  name  = "database_password"
  type  = "String"
  value = random_password.database_password.result
}