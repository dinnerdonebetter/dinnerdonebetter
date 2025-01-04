terraform {
  required_version = "1.10.3"

  backend "remote" {
    organization = "dinnerdonebetter"

    workspaces {
      name = "dev-infra"
    }
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.14.1"
    }
  }
}
