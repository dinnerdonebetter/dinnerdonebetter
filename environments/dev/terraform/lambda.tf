locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 128
  timeout        = 15
}

data "archive_file" "dummy_zip" {
  type        = "zip"
  output_path = "${path.module}/data_changes_lambda.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
}

resource "aws_vpc_endpoint" "ssm_endpoint" {
  vpc_id       = aws_vpc.main.id
  service_name = "com.amazonaws.${local.aws_region}.ssm"
}
