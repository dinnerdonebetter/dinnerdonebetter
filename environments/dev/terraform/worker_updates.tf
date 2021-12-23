resource "aws_sqs_queue" "updates_dead_letter" {
  name                        = "updates_dead_letter"
  sqs_managed_sse_enabled     = true
  content_based_deduplication = true
}

resource "aws_sqs_queue" "updates_queue" {
  name                        = "updates"
  sqs_managed_sse_enabled     = true
  content_based_deduplication = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.updates_dead_letter.arn
    maxReceiveCount     = 1
  })
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

  vpc_config {
    subnet_ids = concat(
      [for x in aws_subnet.public_subnets : x.id],
      [for x in aws_subnet.private_subnets : x.id],
    )
    security_group_ids = [
      aws_security_group.lambda_workers.id,
    ]
  }

  layers = [
    local.collector_layer_arns.us-east-1,
  ]

  filename = data.archive_file.dummy_zip.output_path

  depends_on = [
    aws_cloudwatch_log_group.updates_worker_lambda_logs,
  ]
}

resource "aws_lambda_event_source_mapping" "updates_mapping" {
  event_source_arn = aws_sqs_queue.updates_queue.arn
  function_name    = aws_lambda_function.updates_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "updates_worker_lambda_logs" {
  name              = "/aws/lambda/updates_worker"
  retention_in_days = local.log_retention_period_in_days
}
