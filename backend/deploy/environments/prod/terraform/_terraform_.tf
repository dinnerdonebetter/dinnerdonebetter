terraform {
  required_version = "1.14.5"

  backend "remote" {
    organization = "dinnerdonebetter"

    workspaces {
      name = "prod-backend"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.18.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "6.14.1"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.35.1"
    }
    algolia = {
      source  = "philippe-vandermoere/algolia"
      version = "0.7.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "4.25.0"
    }
  }
}
