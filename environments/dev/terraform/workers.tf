locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 1024
  timeout        = 30
}

data "aws_ecr_repository" "writer_worker" {
  name = "writer_worker"
}

resource "aws_ecr_repository" "updater_worker" {
  name = "updater_worker"
}

resource "aws_ecr_repository" "archiver_worker" {
  name = "archiver_worker"
}

resource "aws_ecr_repository" "data_changes_worker" {
  name = "data_changes_worker"
}

resource "aws_ecr_repository" "chore_worker" {
  name = "chore_worker"
}

resource "aws_lambda_function" "writes_worker_lambda" {
  function_name = "writes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = format("%s:latest", data.aws_ecr_repository.writer_worker.repository_url)
}

resource "aws_lambda_function" "updates_worker_lambda" {
  function_name = "updates_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = format("%s:latest", aws_ecr_repository.updater_worker.repository_url)

  depends_on = [aws_ecr_repository.updater_worker]
}

resource "aws_lambda_function" "archives_worker_lambda" {
  function_name = "archives_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = format("%s:latest", aws_ecr_repository.archiver_worker.repository_url)

  depends_on = [aws_ecr_repository.archiver_worker]
}

resource "aws_lambda_function" "data_changes_worker_lambda" {
  function_name = "data_changes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = format("%s:latest", aws_ecr_repository.data_changes_worker.repository_url)

  depends_on = [aws_ecr_repository.data_changes_worker]
}

resource "aws_lambda_function" "chores_worker_lambda" {
  function_name = "chores_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = format("%s:latest", aws_ecr_repository.chore_worker.repository_url)

  depends_on = [aws_ecr_repository.chore_worker]
}

resource "aws_cloudwatch_event_rule" "every_minute" {
  name                = "every-minute"
  description         = "Fires every minute"
  schedule_expression = "rate(1 minute)"
}

resource "aws_cloudwatch_event_target" "run_chores_every_minute" {
  rule      = aws_cloudwatch_event_rule.every_minute.name
  target_id = "chores_worker"
  arn       = aws_lambda_function.chores_worker_lambda.arn
}