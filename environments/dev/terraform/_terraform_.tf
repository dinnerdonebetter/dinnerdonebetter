terraform {
  required_version = "~> 1.0.11"

  backend "remote" {
    organization = "prixfixe"

    workspaces {
      name = "dev-API"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "~> 1.14"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
    ec = {
      source  = "elastic/ec"
      version = "0.3.0"
    }
    honeycombio = {
      source  = "kvrhdn/honeycombio"
      version = "~> 0.1.0"
    }
  }
}
