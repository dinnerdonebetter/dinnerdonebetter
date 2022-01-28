
data "aws_iam_policy_document" "allow_to_manipulate_queues" {
  statement {
    effect = "Allow"
    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
    ]
    resources = [
      aws_sqs_queue.data_changes_queue.arn,
    ]
  }
}

data "aws_iam_policy_document" "allow_parameter_store_access" {
  statement {
    effect = "Allow"
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:DescribeParameters",
      "ssm:GetParametersByPath",
    ]
    resources = [
      aws_ssm_parameter.service_config.arn,
      aws_ssm_parameter.worker_config.arn,
      aws_ssm_parameter.data_changes_queue_parameter.arn,
      aws_ssm_parameter.database_url.arn,
      aws_ssm_parameter.sendgrid_token.arn,
      aws_ssm_parameter.segment_token.arn,
      aws_ssm_parameter.cookie_hash_key.arn,
      aws_ssm_parameter.cookie_block_key.arn,
      aws_ssm_parameter.paseto_local_key.arn,
    ]
  }
}

data "aws_iam_policy_document" "allow_to_decrypt_parameters" {
  statement {
    effect = "Allow"
    actions = [
      "kms:Decrypt",
      "kms:Verify",
    ]
    resources = [
      aws_kms_key.parameter_store_key.arn,
    ]
  }
}

data "aws_iam_policy_document" "allowed_to_write_traces" {
  statement {
    effect = "Allow"
    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords",
    ]
    resources = [
      "*",
    ]
  }
}

data "aws_iam_policy_document" "allowed_to_network_in_the_vpc" {
  statement {
    effect = "Allow"
    actions = [
      "ec2:CreateNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DeleteNetworkInterface",
    ]
    resources = [
      "*",
    ]
  }
}

data "aws_iam_policy_document" "allowed_to_write_to_the_database" {
  statement {
    effect = "Allow"
    actions = [
      "rds-data:*"
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role" "worker_lambda_role" {
  name = "Worker"

  inline_policy {
    name   = "allow_sqs_queue_access"
    policy = data.aws_iam_policy_document.allow_to_manipulate_queues.json
  }

  inline_policy {
    name   = "allow_ssm_access"
    policy = data.aws_iam_policy_document.allow_parameter_store_access.json
  }

  inline_policy {
    name   = "allow_decrypt_ssm_parameters"
    policy = data.aws_iam_policy_document.allow_to_decrypt_parameters.json
  }

  inline_policy {
    name   = "allowed_to_write_traces"
    policy = data.aws_iam_policy_document.allowed_to_write_traces.json
  }

  inline_policy {
    name   = "allowed_to_network_in_the_vpc"
    policy = data.aws_iam_policy_document.allowed_to_network_in_the_vpc.json
  }

  inline_policy {
    name   = "allowed_to_write_to_the_database"
    policy = data.aws_iam_policy_document.allowed_to_write_to_the_database.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
  ]

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = "sts:AssumeRole",
        Principal = {
          Service = "lambda.amazonaws.com",
        },
      },
    ],
  })
}
