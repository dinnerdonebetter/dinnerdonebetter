resource "aws_sqs_queue" "archives_dead_letter" {
  name = "archives_dead_letter"
}

resource "aws_sqs_queue" "archives_queue" {
  name = "archives"

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.archives_dead_letter.arn
    maxReceiveCount     = 5
  })
}

resource "aws_ssm_parameter" "archives_queue_parameter" {
  name  = "PRIXFIXE_ARCHIVES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.archives_queue.url
}

data "archive_file" "archives_dummy" {
  type        = "zip"
  output_path = "${path.module}/writes_lambda.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
}

resource "aws_lambda_function" "archives_worker_lambda" {
  function_name = "archiver_worker"
  handler       = "archiver_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  tracing_config {
    mode = "Active"
  }

  filename = data.archive_file.archives_dummy.output_path
}

resource "aws_lambda_event_source_mapping" "archives_mapping" {
  event_source_arn = aws_sqs_queue.archives_queue.arn
  function_name    = aws_lambda_function.archives_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "archives_worker_lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.archives_worker_lambda.function_name}"
  retention_in_days = 14
}

resource "aws_cloudwatch_log_subscription_filter" "archives_worker_lambda_subscription_filter" {
  name            = format("%s-postgres-log-group-subscription", aws_lambda_function.archives_worker_lambda.function_name)
  log_group_name  = aws_cloudwatch_log_group.archives_worker_lambda_logs.name
  filter_pattern  = ""
  destination_arn = data.aws_lambda_function.cloudwatch_logs.arn
}