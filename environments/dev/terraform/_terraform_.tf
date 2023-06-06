terraform {
  required_version = "1.4.6"

  backend "remote" {
    organization = "dinnerdonebetter"

    workspaces {
      name = "dev-backend"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.3.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "4.43.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.2.0"
    }
    algolia = {
      source = "philippe-vandermoere/algolia"
      version = "0.7.0"
    }
  }
}
