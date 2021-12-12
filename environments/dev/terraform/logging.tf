data "aws_cloudwatch_log_group" "dev" {
  name = "development"
}

data "aws_lambda_function" "cloudwatch_logs" {
  function_name = "honeycomb-cloudwatch-logs-integration"
}

resource "aws_cloudwatch_log_subscription_filter" "cloudwatch_subscription_filter" {
  name            = format("%s-log-group-subscription", data.aws_cloudwatch_log_group.dev.name)
  log_group_name  = data.aws_cloudwatch_log_group.dev.name
  filter_pattern  = ""
  destination_arn = data.aws_lambda_function.cloudwatch_logs.arn
}
