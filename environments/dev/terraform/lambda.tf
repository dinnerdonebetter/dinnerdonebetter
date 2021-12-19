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

resource "aws_security_group" "lambda_workers" {
  name        = "lambda"
  description = "Lambda group"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_vpc_endpoint" "ssm_endpoint" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${local.aws_region}.ssm"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.lambda_workers.id,
  ]

  private_dns_enabled = true
}


resource "aws_vpc_endpoint" "kms_endpoint" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${local.aws_region}.kms"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.lambda_workers.id,
  ]

  private_dns_enabled = true
}
