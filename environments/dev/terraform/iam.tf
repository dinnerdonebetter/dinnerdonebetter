
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
      aws_sqs_queue.writes_queue.arn,
      aws_sqs_queue.updates_queue.arn,
      aws_sqs_queue.archives_queue.arn,
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
      aws_ssm_parameter.opentelemetry_collector_config.arn,
      aws_ssm_parameter.writes_queue_parameter.arn,
      aws_ssm_parameter.updates_queue_parameter.arn,
      aws_ssm_parameter.archives_queue_parameter.arn,
      aws_ssm_parameter.data_changes_queue_parameter.arn,
      aws_ssm_parameter.database_url.arn,
      aws_ssm_parameter.sendgrid_token.arn,
      aws_ssm_parameter.segment_token.arn,
      aws_ssm_parameter.cookie_hash_key.arn,
      aws_ssm_parameter.cookie_block_key.arn,
      aws_ssm_parameter.paseto_local_key.arn,
      aws_ssm_parameter.search_url.arn,
      aws_ssm_parameter.search_username.arn,
      aws_ssm_parameter.search_password.arn,
      aws_ssm_parameter.pubsub_server_url.arn,
      aws_ssm_parameter.pubsub_server_username.arn,
      aws_ssm_parameter.pubsub_server_password.arn,
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
      "ec2:DescribeNetworkInterfaces",
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

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
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

resource "aws_iam_role" "server_lambda_role" {
  name = "APIServer"

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

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com",
        },
      },
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com",
        },
      },
    ],
  })
}