resource "aws_sqs_queue" "data_changes_dead_letter" {
  name                    = "data_changes_dead_letter"
  sqs_managed_sse_enabled = true
}

resource "aws_sqs_queue" "data_changes_queue" {
  name = "data_changes"
  sqs_managed_sse_enabled = true
}

resource "aws_ssm_parameter" "data_changes_queue_parameter" {
  name  = "PRIXFIXE_DATA_CHANGES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.data_changes_queue.arn
}

resource "aws_lambda_function" "data_changes_worker_lambda" {
  function_name = "data_changes_worker"
  handler       = "data_changes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
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
    aws_cloudwatch_log_group.data_changes_worker_lambda_logs,
  ]
}

resource "aws_lambda_event_source_mapping" "data_changes_mapping" {
  event_source_arn = aws_sqs_queue.data_changes_queue.arn
  function_name    = aws_lambda_function.data_changes_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "data_changes_worker_lambda_logs" {
  name              = "/aws/lambda/data_changes"
  retention_in_days = local.log_retention_period_in_days
}
