resource "aws_sqs_queue" "writes_queue" {
  name       = "writes.fifo"
  fifo_queue = true
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

  filename         = "writer_lambda.zip"
  source_code_hash = filebase64sha256("writer_lambda.zip")
}

resource "aws_lambda_event_source_mapping" "writes_mapping" {
  event_source_arn = aws_sqs_queue.writes_queue.arn
  function_name    = aws_lambda_function.writes_worker_lambda.arn
}