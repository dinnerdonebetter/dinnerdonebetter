locals {
  database_username = "prixfixe_api"
  database_name     = "prixfixe"
  cluster_name      = "api-database"
}

resource "aws_cloudwatch_log_group" "database_logs" {
  name              = "/aws/rds/cluster/${local.cluster_name}/postgresql"
  retention_in_days = local.log_retention_period_in_days
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
  cluster_identifier = local.cluster_name
  database_name      = local.database_name
  engine             = "aurora-postgresql"

  engine_mode = "serverless"
  scaling_configuration {
    auto_pause               = true
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

resource "aws_security_group" "database" {
  name        = "postgres"
  description = "Allow Postgres traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description      = "Postgres from VPC"
    from_port        = 5432
    to_port          = 5432
    protocol         = "tcp"
    cidr_blocks      = [aws_vpc.main.cidr_block]
    ipv6_cidr_blocks = [aws_vpc.main.ipv6_cidr_block]
    security_groups = [
      aws_security_group.api_service.id,
      aws_security_group.load_balancer.id,
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "honeycombio_dataset" "dev_postgres_rds_logs" {
  name = "dev_postgres_rds_logs"
}

resource "aws_lambda_function" "honeycomb_postgres_rds_logs" {
  # change me to your region
  s3_bucket     = "honeycomb-integrations-us-east-1"
  s3_key        = "agentless-integrations-for-aws/LATEST/ingest-handlers.zip"
  function_name = "honeycomb-postgres-rds-logs"
  role          = aws_iam_role.honeycomb_logs.arn
  handler       = "postgresql-handler"
  runtime       = "go1.x"
  memory_size   = "128"

  environment {
    variables = {
      ENVIRONMENT         = "dev"
      HONEYCOMB_WRITE_KEY = var.HONEYCOMB_API_KEY
      DATASET             = honeycombio_dataset.dev_postgres_rds_logs.name
      SAMPLE_RATE         = "1"
      SCRUB_QUERY         = "false"
      HONEYCOMB_DEBUG     = true
    }
  }
}

resource "aws_lambda_permission" "allow_database_logs_from_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.honeycomb_postgres_rds_logs.arn
  principal     = "logs.amazonaws.com"
}

resource "aws_cloudwatch_log_subscription_filter" "database_logs_subscription_filter" {
  name            = "honeycomb-postgres-rds-subscription"
  log_group_name  = aws_cloudwatch_log_group.database_logs.name
  filter_pattern  = ""
  destination_arn = aws_lambda_function.honeycomb_postgres_rds_logs.arn
}
