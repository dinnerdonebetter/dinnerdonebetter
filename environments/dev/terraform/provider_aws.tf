variable "AWS_ACCESS_KEY" {}
variable "AWS_SECRET_ACCESS_KEY" {}

provider "aws" {
  region     = local.aws_region
  access_key = var.AWS_ACCESS_KEY
  secret_key = var.AWS_SECRET_ACCESS_KEY

  default_tags {
    tags = {
      Environment = "dev"
      CreatedVia  = "terraform"
    }
  }
}

data "aws_caller_identity" "current" {}

output "caller_id" {
  value = data.aws_caller_identity.current.id
}

output "caller_account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "user_id" {
  value = data.aws_caller_identity.current.account_id
}