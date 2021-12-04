terraform {
  required_version = "~> 1.0.11"

  backend "remote" {
    organization = "prixfixe"

    workspaces {
      name = "dev-API"
    }
  }
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}
