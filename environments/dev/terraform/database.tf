module "aurora_postgresql" {
  source = "terraform-aws-modules/rds-aurora/aws"

  name              = "dev-api-database"
  engine            = "aurora-postgresql"
  engine_mode       = "serverless"
  storage_encrypted = true

  # create_security_group = true

  monitoring_interval = 60

  apply_immediately   = true
  skip_final_snapshot = true

  scaling_configuration = {
    auto_pause               = true
    min_capacity             = 1
    max_capacity             = 1
    seconds_until_auto_pause = 300
    timeout_action           = "ForceApplyCapacityChange"
  }

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}