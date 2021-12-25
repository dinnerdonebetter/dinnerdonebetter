variable "AWS_ACCESS_KEY_ID" {}
variable "AWS_SECRET_ACCESS_KEY" {}

provider "aws" {
  region     = local.aws_region
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY

  default_tags {
    tags = {
      Environment = "dev"
      CreatedVia  = "terraform"
    }
  }
}

data "aws_caller_identity" "current" {}

locals {
  aws_region = "us-east-1"

  log_retention_period_in_days = 14
}