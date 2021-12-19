resource "aws_sqs_queue" "writes_dead_letter" {
  name                    = "writes_dead_letter"
  sqs_managed_sse_enabled = true
}

resource "aws_sqs_queue" "writes_queue" {
  name = "writes"
  # sqs_managed_sse_enabled = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.writes_dead_letter.arn
    maxReceiveCount     = 1
  })
}

resource "aws_ssm_parameter" "writes_queue_parameter" {
  name  = "PRIXFIXE_WRITES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.writes_queue.url
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

  #  layers = [
  #    local.collector_layer_arns.us-east-1,
  #  ]

  filename = data.archive_file.dummy_zip.output_path

  depends_on = [
    aws_cloudwatch_log_group.writes_worker_lambda_logs,
  ]
}

resource "aws_lambda_event_source_mapping" "writes_mapping" {
  event_source_arn = aws_sqs_queue.writes_queue.arn
  function_name    = aws_lambda_function.writes_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "writes_worker_lambda_logs" {
  name              = "/aws/lambda/writes_worker"
  retention_in_days = local.log_retention_period_in_days
}
