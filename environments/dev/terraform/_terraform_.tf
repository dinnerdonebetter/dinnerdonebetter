terraform {
  required_version = "~> 1.1.2"

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
    google = {
      source = "hashicorp/google"
      version = "4.8.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }
}
