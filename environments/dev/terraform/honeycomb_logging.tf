# This example should help you get going running the generic JSON integration for Cloudwatch ogs
# It won't work out of the box - you will need to update some environment variables, and possibly tweak
# the configuration to work within your TF environment.
resource "aws_iam_role" "honeycomb_logs" {
  name = "honeycomb-cloudwatch-logs-lambda-role"

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

resource "aws_iam_role_policy" "lambda_log_policy" {
  name   = "lambda-logs-policy"
  role   = aws_iam_role.honeycomb_logs.id
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}
