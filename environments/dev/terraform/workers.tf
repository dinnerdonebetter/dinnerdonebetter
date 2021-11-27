locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 1024
  timeout        = 30
}

data "aws_iam_role" "worker_lambda_role" {
  name = "Workers"
}

resource "aws_lambda_function" "writes_worker_lambda" {
  function_name = "writes_worker"
  role          = data.aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  filename         = "writer_lambda.zip"
  source_code_hash = filebase64sha256("writer_lambda.zip")


  tags = merge(var.default_tags, {})
}

resource "aws_lambda_function" "updates_worker_lambda" {
  function_name = "updates_worker"
  role          = data.aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  filename         = "updater_lambda.zip"
  source_code_hash = filebase64sha256("updater_lambda.zip")

  tags = merge(var.default_tags, {})
}

resource "aws_lambda_function" "archives_worker_lambda" {
  function_name = "archives_worker"
  role          = data.aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  filename         = "archiver_lambda.zip"
  source_code_hash = filebase64sha256("archiver_lambda.zip")

  tags = merge(var.default_tags, {})
}

resource "aws_lambda_function" "data_changes_worker_lambda" {
  function_name = "data_changes_worker"
  role          = data.aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  filename         = "data_changes_lambda.zip"
  source_code_hash = filebase64sha256("data_changes_lambda.zip")

  tags = merge(var.default_tags, {})
}

resource "aws_lambda_function" "chores_worker_lambda" {
  function_name = "chores_worker"
  role          = data.aws_iam_role.worker_lambda_role.arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = local.timeout

  filename         = "chores_lambda.zip"
  source_code_hash = filebase64sha256("chores_lambda.zip")

  tags = merge(var.default_tags, {})
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