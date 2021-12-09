locals {
  database_username = "prixfixe_api"
  database_name     = "prixfixe"
}

resource "random_password" "database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "aws_db_subnet_group" "db_subnet" {
  name        = "dev"
  description = "dev environment database subnet group"
  subnet_ids  = [for x in aws_subnet.private_subnets : x.id]
}

resource "aws_rds_cluster" "api_database" {
  cluster_identifier = "api-database"
  database_name      = local.database_name
  engine             = "aurora-postgresql"

  engine_mode = "serverless"
  scaling_configuration {
    # auto_p//use               = true
    min_capacity             = 2
    max_capacity             = 2
    seconds_until_auto_pause = 300
    timeout_action           = "ForceApplyCapacityChange"
  }

  master_username              = local.database_username
  master_password              = random_password.database_password.result
  preferred_backup_window      = "05:00-08:00"
  preferred_maintenance_window = "sat:01:00-sat:04:00"
  apply_immediately            = true
  enable_http_endpoint         = true
  storage_encrypted            = true
  skip_final_snapshot          = true
  copy_tags_to_snapshot        = true
  backup_retention_period      = 7

  db_subnet_group_name = aws_db_subnet_group.db_subnet.name
  vpc_security_group_ids = [
    aws_security_group.database.id,
  ]
}

resource "aws_ssm_parameter" "database_url" {
  name = "PRIXFIXE_DATABASE_CONNECTION_STRING"
  type = "String"
  value = format(
    "user=%s dbname=%s password='%s' host=%s port=%s",
    local.database_username,
    local.database_name,
    random_password.database_password.result,
    aws_rds_cluster.api_database.endpoint,
    aws_rds_cluster.api_database.port,
  )
}


resource "aws_secretsmanager_secret" "dev_database" {
  name = format("rds-db-credentials/%s/%s", aws_rds_cluster.api_database.cluster_resource_id, local.database_username)
}

resource "aws_secretsmanager_secret_version" "dev_database" {
  secret_id = aws_secretsmanager_secret.dev_database.id
  secret_string = jsonencode({
    dbInstanceIdentifier = "api-database",
    engine               = aws_rds_cluster.api_database.engine,
    host                 = aws_rds_cluster.api_database.endpoint,
    port                 = aws_rds_cluster.api_database.port,
    resourceId           = aws_rds_cluster.api_database.cluster_resource_id,
    username             = local.database_username,
    password             = random_password.database_password.result
  })
}