variable "GRAFANA_CLOUD_API_KEY" {}

locals {
  grafana_stack = "dinnerdonebetter"
}

provider "grafana" {
  url                       = "https://${local.grafana_stack}.grafana.net/"
  cloud_access_policy_token = var.GRAFANA_CLOUD_API_KEY
}

# resource "grafana_cloud_plugin_installation" "github" {
#   stack_slug = local.grafana_stack
#   slug       = "grafana-github-datasource"
#   version    = "2.0.1"
# }
