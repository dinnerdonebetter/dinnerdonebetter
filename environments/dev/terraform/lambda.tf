locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 128
  timeout        = 15

  cloudwatch_exclude_lambda_events_filter = "-\"START RequestId\" -\"END RequestId\" -\"REPORT RequestId\""
}

resource "aws_cloudwatch_log_group" "honeycomb_worker_logs" {
  name              = "/aws/lambda/honeycomb-worker-logs-integration"
  retention_in_days = local.log_retention_period_in_days
}

resource "aws_lambda_function" "worker_log_sync" {
  function_name = "honeycomb-worker-logs-integration"
  s3_bucket     = "honeycomb-integrations-us-east-1"
  s3_key        = "agentless-integrations-for-aws/LATEST/ingest-handlers.zip"
  handler       = "cloudwatch-handler"
  runtime       = "go1.x"
  memory_size   = "128"
  role          = aws_iam_role.honeycomb_logs.arn

  environment {
    variables = {
      ENVIRONMENT         = "dev"
      PARSER_TYPE         = "json"
      HONEYCOMB_WRITE_KEY = var.HONEYCOMB_API_KEY
      DATASET             = "dev_api_server_logs"
      SAMPLE_RATE         = "1"
      TIME_FIELD_NAME     = "time"
      TIME_FIELD_FORMAT   = "%s(%L)?"
    }
  }
}

resource "aws_lambda_permission" "allow_worker_logs_from_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.worker_log_sync.arn
  principal     = "logs.amazonaws.com"
}
