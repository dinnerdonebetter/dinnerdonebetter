resource "aws_prometheus_workspace" "dev" {
  alias = "dev"
}

resource "aws_ssm_parameter" "prometheus_endpoint" {
  name  = "PRIXFIXE_PROMETHEUS_ENDPOINT"
  type  = "String"
  value = aws_prometheus_workspace.dev.prometheus_endpoint
}

resource "aws_ssm_parameter" "opentelemetry_collector_config" {
  name  = "otel-collector-config"
  type  = "String"
  value = file("${path.module}/opentelemetry-config.yaml")
}

data "aws_iam_policy_document" "opentelemetry_collector_policy" {
  version = "2012-10-17"

  statement {
    effect = "Allow"
    actions = [
      "logs:PutLogEvents",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogStreams",
      "logs:DescribeLogGroups",
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords",
      "xray:GetSamplingRules",
      "xray:GetSamplingTargets",
      "xray:GetSamplingStatisticSummaries",
      "cloudwatch:PutMetricData",
      "ec2:DescribeVolumes",
      "ec2:DescribeTags",
      "ssm:GetParameters",
      "aps:RemoteWrite",
      "aps:GetSeries",
      "aps:GetLabels",
      "aps:GetMetricMetadata",
    ]
    resources = ["*"]
  }
}
