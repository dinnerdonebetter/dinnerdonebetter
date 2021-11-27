resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_BASE_CONFIG"
  type  = "String"
  value =  "${file("../config_files/service-config.json")}"
}