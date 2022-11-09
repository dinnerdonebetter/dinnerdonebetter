terraform {
  required_version = "~> 1.3.4"

  backend "remote" {
    organization = "prixfixe"

    workspaces {
      name = "prixfixe-backend-dev"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "3.9.1"
    }
    google = {
      source  = "hashicorp/google"
      version = "4.11.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }
}
