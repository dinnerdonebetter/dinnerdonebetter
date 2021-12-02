variable "AWS_ACCESS_KEY" {}
variable "AWS_SECRET_ACCESS_KEY" {}

provider "aws" {
  region     = "us-east-1"
  access_key = var.AWS_ACCESS_KEY
  secret_key = var.AWS_SECRET_ACCESS_KEY

  default_tags {
    tags = {
      Environment = "dev"
      CreatedVia  = "terraform"
      CreatedAt   = timestamp()
    }
  }
}
