variable "ELASTIC_CLOUD_API_KEY" {}

provider "ec" {
  apikey = var.ELASTIC_CLOUD_API_KEY
}
