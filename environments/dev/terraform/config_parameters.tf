resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_BASE_API_SERVER_CONFIG"
  type  = "String"
  value = file("${path.module}/service-config.json")
}

resource "aws_ssm_parameter" "worker_config" {
  name  = "PRIXFIXE_BASE_WORKER_CONFIG"
  type  = "String"
  value = file("${path.module}/worker-config.json")
}
