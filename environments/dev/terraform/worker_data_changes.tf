resource "aws_sqs_queue" "data_changes_dead_letter" {
  name = "data_changes_dead_letter"
}

resource "aws_sns_topic" "data_changes_queue" {
  name = "data_changes"

  #  redrive_policy = jsonencode({
  #    deadLetterTargetArn = aws_sqs_queue.data_changes_dead_letter.arn
  #    maxReceiveCount     = 5
  #  })
}

resource "aws_ssm_parameter" "data_changes_queue_parameter" {
  name  = "PRIXFIXE_DATA_CHANGES_QUEUE_URL"
  type  = "String"
  value = aws_sns_topic.data_changes_queue.arn
}

data "archive_file" "data_changes_dummy" {
  type        = "zip"
  output_path = "${path.module}/data_changes_lambda.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
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

  layers = [
    local.collector_layer_arns.us-east-1,
  ]

  filename = data.archive_file.data_changes_dummy.output_path
}

resource "aws_sns_topic_subscription" "data_changes_mapping" {
  topic_arn = aws_sns_topic.data_changes_queue.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.data_changes_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "data_changes_worker_lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.data_changes_worker_lambda.function_name}"
  retention_in_days = local.log_retention_period_in_days
}
