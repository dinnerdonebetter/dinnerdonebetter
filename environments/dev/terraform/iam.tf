resource "aws_iam_role" "worker_lambda_role" {
  name = "Worker"

  inline_policy {
    name = "allow_sqs_queue_access"

    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Action = [
            "sqs:*",
          ],
          Effect = "Allow",
          Resource = [
            aws_sqs_queue.writes_queue.arn,
            aws_sqs_queue.updates_queue.arn,
            aws_sqs_queue.archives_queue.arn,
            aws_sqs_queue.data_changes_queue.arn,
            aws_sqs_queue.chores_queue.arn,
          ],
        }
      ]
    })
  }

  inline_policy {
    name = "allow ssm access"
    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Action = [
            "ssm:GetParameter",
          ],
          Effect = "Allow",
          Resource = [
            aws_ssm_parameter.service_config.arn,
            aws_ssm_parameter.writes_queue_parameter.arn,
            aws_ssm_parameter.updates_queue_parameter.arn,
            aws_ssm_parameter.archives_queue_parameter.arn,
            aws_ssm_parameter.data_changes_queue_parameter.arn,
            aws_ssm_parameter.chores_queue_parameter.arn,
            aws_ssm_parameter.database_url.arn,
          ],
        }
      ]
    })
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
    name = "allow_sqs_queue_access"

    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Action = [
            "sqs:*",
          ],
          Effect   = "Allow",
          Resource = "*",
        },
      ],
    })
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