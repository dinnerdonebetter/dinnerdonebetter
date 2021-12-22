resource "aws_sqs_queue" "archives_dead_letter" {
  name                    = "archives_dead_letter"
  sqs_managed_sse_enabled = true
}

resource "aws_sqs_queue" "archives_queue" {
  name                    = "archives"
  sqs_managed_sse_enabled = true

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

  vpc_config {
    subnet_ids = concat(
      [for x in aws_subnet.public_subnets : x.id],
      [for x in aws_subnet.private_subnets : x.id],
    )
    security_group_ids = [
      aws_security_group.lambda_workers.id,
    ]
  }

  #  layers = [
  #    local.collector_layer_arns.us-east-1,
  #  ]

  filename = data.archive_file.dummy_zip.output_path

  depends_on = [
    aws_cloudwatch_log_group.archives_worker_lambda_logs,
  ]
}

resource "aws_lambda_event_source_mapping" "archives_mapping" {
  event_source_arn = aws_sqs_queue.archives_queue.arn
  function_name    = aws_lambda_function.archives_worker_lambda.arn
}

resource "aws_cloudwatch_log_group" "archives_worker_lambda_logs" {
  name              = "/aws/lambda/archives_worker"
  retention_in_days = local.log_retention_period_in_days
}
