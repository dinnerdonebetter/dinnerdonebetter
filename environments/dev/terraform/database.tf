locals {
  database_username = "prixfixe_api"
  database_name     = "prixfixe"
  cluster_name      = "api-database"
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

resource "aws_db_instance" "api_database" {
  allocated_storage = 10
  engine            = "postgres"
  engine_version    = "12"
  instance_class    = "db.t2.micro"
  name              = local.database_name

  username = local.database_username
  password = random_password.database_password.result

  skip_final_snapshot = true
  port                = 5423
  # storage_encrypted = true # InvalidParameterCombination: DB Instance class db.t2.micro does not support encryption at rest

  backup_window      = "05:00-08:00"
  maintenance_window = "sat:01:00-sat:04:00"

  db_subnet_group_name = aws_db_subnet_group.db_subnet.name
  vpc_security_group_ids = [
    aws_security_group.database.id,
  ]

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
}

resource "aws_ssm_parameter" "database_url" {
  name = "PRIXFIXE_DATABASE_CONNECTION_STRING"
  type = "String"
  value = format(
    "user=%s dbname=%s password='%s' host=%s port=%s",
    local.database_username,
    local.database_name,
    random_password.database_password.result,
    aws_db_instance.api_database.endpoint,
    aws_db_instance.api_database.port,
  )
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
      aws_security_group.lambda_workers.id,
    ]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}
