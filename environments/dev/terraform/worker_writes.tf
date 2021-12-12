resource "aws_sqs_queue" "writes_dead_letter" {
  name = "writes_dead_letter"
}

resource "aws_sqs_queue" "writes_queue" {
  name = "writes"

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.writes_dead_letter.arn
    maxReceiveCount     = 5
  })
}

resource "aws_ssm_parameter" "writes_queue_parameter" {
  name  = "PRIXFIXE_WRITES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.writes_queue.url
}

data "archive_file" "writes_lambda_dummy" {
  type        = "zip"
  output_path = "${path.module}/writes_lambda.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
}

resource "aws_lambda_function" "writes_worker_lambda" {
  function_name = "writes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler       = "writes_worker"
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  tracing_config {
    mode = "Active"
  }

  filename = data.archive_file.writes_lambda_dummy.output_path
}

resource "aws_lambda_event_source_mapping" "writes_mapping" {
  event_source_arn = aws_sqs_queue.writes_queue.arn
  function_name    = aws_lambda_function.writes_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "writes_worker_lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.writes_worker_lambda.function_name}"
  retention_in_days = local.log_retention_period_in_days
}

#resource "aws_cloudwatch_log_subscription_filter" "writes_worker_lambda_subscription_filter" {
#  name            = format("%s-log-group-subscription", aws_lambda_function.writes_worker_lambda.function_name)
#  log_group_name  = aws_cloudwatch_log_group.writes_worker_lambda_logs.name
#  filter_pattern  = local.cloudwatch_exclude_lambda_events_filter
#  destination_arn = aws_lambda_function.cloudwatch_logs.arn
#}