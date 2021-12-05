resource "aws_sqs_queue" "updates_queue" {
  name       = "updates.fifo"
  fifo_queue = true
}

resource "aws_ssm_parameter" "updates_queue_parameter" {
  name  = "PRIXFIXE_UPDATES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.updates_queue.url
}

resource "aws_lambda_function" "updates_worker_lambda" {
  function_name = "updates_worker"
  handler       = "updates_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  tracing_config {
    mode = "Active"
  }

  #   filename = "updates_lambda.zip"
}

resource "aws_lambda_event_source_mapping" "updates_mapping" {
  event_source_arn = aws_sqs_queue.updates_queue.arn
  function_name    = aws_lambda_function.updates_worker_lambda.arn
}