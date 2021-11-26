locals {
  lambda_runtime = "go1.x"
  lambda_handler  = "main"
  memory_size = 1024
  timeout = 30
}

resource "aws_iam_role" "worker_lambda_role" {
  name = "worker_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "writes_worker_lambda" {
  function_name = "writes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler = local.lambda_handler
  runtime = local.lambda_runtime
  memory_size       = local.memory_size
  timeout           = local.timeout

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_function" "updates_worker_lambda" {
  function_name = "updates_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler = local.lambda_handler
  runtime = local.lambda_runtime
  memory_size       = local.memory_size
  timeout           = local.timeout

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_function" "archives_worker_lambda" {
  function_name = "archives_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler = local.lambda_handler
  runtime = local.lambda_runtime
  memory_size       = local.memory_size
  timeout           = local.timeout

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_function" "data_changes_worker_lambda" {
  function_name = "data_changes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  handler = local.lambda_handler
  runtime = local.lambda_runtime
  memory_size       = local.memory_size
  timeout           = local.timeout

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}