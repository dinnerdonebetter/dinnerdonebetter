terraform {
  required_version = "1.10.3"

  backend "remote" {
    organization = "dinnerdonebetter"

    workspaces {
      name = "dev-infra"
    }
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.40.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "6.14.1"
    }
    sendgrid = {
      source  = "Trois-Six/sendgrid"
      version = "0.2.1"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "3.16.0"
    }
  }
}
