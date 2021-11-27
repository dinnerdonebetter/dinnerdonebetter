data "local_file" "service_config" {
  filename = "${path.module}/../config_files/service-config.json"
}

resource "aws_ssm_parameter" "service_config" {
  name  = "PRIXFIXE_BASE_CONFIG"
  type  = "String"
  value = data.local_file.service_config.content
}