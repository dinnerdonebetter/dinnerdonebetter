resource "aws_cloudwatch_log_group" "dev" {
  name              = "development"
  retention_in_days = 7
}

data "aws_lambda_function" "cloudwatch_logs" {
  function_name = "honeycomb-cloudwatch-logs-integration"
  s3_bucket     = "honeycomb-integrations-us-east-1"
  s3_key        = "agentless-integrations-for-aws/LATEST/ingest-handlers.zip"
}

resource "aws_cloudwatch_log_subscription_filter" "cloudwatch_subscription_filter" {
  name            = format("%s-log-group-subscription", aws_cloudwatch_log_group.dev.name)
  log_group_name  = aws_cloudwatch_log_group.dev.name
  filter_pattern  = ""
  destination_arn = data.aws_lambda_function.cloudwatch_logs.arn
}
