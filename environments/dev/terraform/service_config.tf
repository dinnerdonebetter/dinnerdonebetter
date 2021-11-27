data "local_file" "service_config" {
  filename = "../config_files/service-config.json"
}

resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_BASE_CONFIG"
  type  = "String"
  value = data.local_file.service_config.content
}