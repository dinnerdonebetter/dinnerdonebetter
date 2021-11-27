resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_BASE_CONFIG"
  type  = "String"
  value = file("${ path.module }/service-config.json")
}