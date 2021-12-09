
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
    resources = [for x in aws_ssm_parameter : x.arn]
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