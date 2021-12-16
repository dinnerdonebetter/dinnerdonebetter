resource "aws_prometheus_workspace" "dev" {
  alias = "dev"
}

resource "aws_ssm_parameter" "prometheus_endpoint" {
  name  = "PRIXFIXE_PROMETHEUS_ENDPOINT"
  type  = "String"
  value = aws_prometheus_workspace.dev.prometheus_endpoint
}
