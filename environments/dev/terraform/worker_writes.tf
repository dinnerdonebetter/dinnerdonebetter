resource "aws_sqs_queue" "writes_queue" {
  name       = "writes.fifo"
  fifo_queue = true
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