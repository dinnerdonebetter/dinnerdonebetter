resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_BASE_CONFIG"
  type  = "String"
  value = file("${path.module}/service-config.json")
}

resource "aws_ssm_parameter" "opentelemetry_collector_config" {
  name  = "otel-collector-config"
  type  = "String"
  value = file("${path.module}/opentelemetry-config.yaml")
}
