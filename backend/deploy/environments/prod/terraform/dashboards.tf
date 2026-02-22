resource "grafana_folder" "dashboards" {
  title = "DinnerDoneBetter Dashboards"
}

resource "grafana_dashboard" "api_server" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/api_server.json")
}

resource "grafana_dashboard" "admin_webapp" {
  folder = grafana_folder.dashboards.id

  config_json = file("${path.module}/dashboards/admin_webapp.json")
}
