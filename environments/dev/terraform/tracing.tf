resource "aws_xray_group" "dev" {
  group_name        = "dev"
  filter_expression = "responsetime > 5"
}

resource "aws_xray_sampling_rule" "trace_all" {
  rule_name      = "example"
  priority       = 1
  version        = 1
  reservoir_size = 1
  fixed_rate     = 1
  url_path       = "*"
  host           = "*"
  http_method    = "*"
  service_type   = "*"
  service_name   = "*"
  resource_arn   = "*"
}