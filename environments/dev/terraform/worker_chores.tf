resource "aws_sqs_queue" "chores_queue" {
  name       = "chores.fifo"
  fifo_queue = true
}

resource "aws_ssm_parameter" "chores_queue_parameter" {
  name  = "PRIXFIXE_CHORES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.chores_queue.url
}

data "archive_file" "chores_dummy" {
  type        = "zip"
  output_path = "${path.module}/archives_lambda.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
}

resource "aws_lambda_function" "chores_worker_lambda" {
  function_name = "chores_worker"
  handler       = "chores_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  tracing_config {
    mode = "Active"
  }

  filename = data.archive_file.chores_dummy.output_path
}

resource "aws_lambda_event_source_mapping" "chores_mapping" {
  event_source_arn = aws_sqs_queue.chores_queue.arn
  function_name    = aws_lambda_function.chores_worker_lambda.arn
}