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
          Effect   = "Allow",
          Resource = "*",
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