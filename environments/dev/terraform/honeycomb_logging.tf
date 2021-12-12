# This example should help you get going running the generic JSON integration for Cloudwatch ogs
# It won't work out of the box - you will need to update some environment variables, and possibly tweak
# the configuration to work within your TF environment.
resource "aws_iam_role" "honeycomb_cloudwatch_logs" {
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
  role   = aws_iam_role.honeycomb_cloudwatch_logs.id
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

# resource "aws_kms_key" "api_secret_encryption" {
#   description = "to encrypt parameter store secrets"
#   key_usage   = "ENCRYPT_DECRYPT"
# }

# resource "aws_iam_role_policy" "lambda_kms_policy" {
#   name   = "lambda-kms-policy"
#   role   = "${aws_iam_role.honeycomb_cloudwatch_logs.id}"
#   policy = <<EOF
# {
#   "Version": "2012-10-17",
#   "Statement": [
#     {
#       "Action": "kms:Decrypt",
#       "Effect": "Allow",
#       "Resource": "${aws_kms_key.api_secret_encryption.arn}"
#     }
#   ]
# }
# EOF
# }

resource "aws_cloudwatch_log_group" "honeycomb_cloudwatch_logs" {
  name              = "/aws/lambda/honeycomb-cloudwatch-logs-integration"
  retention_in_days = local.log_retention_period_in_days
}


resource "aws_lambda_function" "cloudwatch_logs" {
  function_name = "honeycomb-cloudwatch-logs-integration"
  s3_bucket     = "honeycomb-integrations-us-east-1"
  s3_key        = "agentless-integrations-for-aws/LATEST/ingest-handlers.zip"
  handler       = "cloudwatch-handler"
  runtime       = "go1.x"
  memory_size   = "128"
  role          = aws_iam_role.honeycomb_cloudwatch_logs.arn

  environment {
    variables = {
      ENVIRONMENT = "dev"
      PARSER_TYPE = "json"
      # Change this to your encrypted Honeycomb write key or your raw write key (not recommended)
      HONEYCOMB_WRITE_KEY = var.HONEYCOMB_API_KEY
      # # If the write key is encrypted, specify the KMS Key ID used to encrypt your write key
      # # see https://github.com/honeycombio/agentless-integrations-for-aws#encrypting-your-write-key
      # KMS_KEY_ID =  aws_kms_key.api_secret_encryption.id
      DATASET           = "logs"
      SAMPLE_RATE       = "1"
      TIME_FIELD_NAME   = "time"
      TIME_FIELD_FORMAT = "%s(%L)?"
    }
  }
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cloudwatch_logs.arn
  principal     = "logs.amazonaws.com"
}

resource "aws_ssm_parameter" "honeycomb_api_key" {
  name  = "PRIXFIXE_HONEYCOMB_API_KEY"
  type  = "String"
  value = var.HONEYCOMB_API_KEY
}