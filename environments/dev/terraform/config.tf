resource "kubernetes_secret_v1" "config_file" {
  metadata {
    name = "prixfixe-api-configuration"
  }

  data = {
    "service-config.json" = file("${path.module}/service-config.json")
  }
}
