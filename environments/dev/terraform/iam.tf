resource "aws_iam_role" "worker_lambda_role" {
  name = "Worker"

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "ec2.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      },
      {
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      },
      {
        "Effect" : "Allow",
        "Resource" : aws_sqs_queue.writes_queue.arn,
        "Action" : "sts:AssumeRole"
      },
      {
        "Effect" : "Allow",
        "Resource" : aws_sqs_queue.updates_queue.arn,
        "Action" : "sts:AssumeRole"
      },
      {
        "Effect" : "Allow",
        "Resource" : aws_sqs_queue.archives_queue.arn,
        "Action" : "sts:AssumeRole"
      },
      {
        "Effect" : "Allow",
        "Resource" : aws_sqs_queue.data_changes_queue.arn,
        "Action" : "sts:AssumeRole"
      },
    ]
  })
}

resource "aws_iam_role" "server_lambda_role" {
  name = "APIServer"

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "ec2.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      },
      {
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      }
    ]
  })
}