locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 128
  timeout        = 8
}

data "archive_file" "dummy_zip" {
  type        = "zip"
  output_path = "${path.module}/data_changes_lambda.zip"

  source {
    content  = "hello"
    filename = "dummy.txt"
  }
}

resource "aws_security_group" "vpc_endpoints" {
  name        = "vpc_endpoints"
  description = "AWS VPC endpoints"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port        = 443
    to_port          = 443
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = 6379
    to_port          = 6379
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_security_group" "lambda_workers" {
  name        = "workers"
  description = "Lambda group"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_vpc_endpoint" "ssm_endpoint" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${local.aws_region}.ssm"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.vpc_endpoints.id,
  ]

  private_dns_enabled = true
}

resource "aws_vpc_endpoint_subnet_association" "private_ssm_association" {
  for_each = aws_subnet.private_subnets

  subnet_id       = each.value.id
  vpc_endpoint_id = aws_vpc_endpoint.ssm_endpoint.id
}

resource "aws_vpc_endpoint" "kms_endpoint" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${local.aws_region}.kms"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.vpc_endpoints.id,
  ]

  private_dns_enabled = true
}

resource "aws_vpc_endpoint_subnet_association" "private_kms_association" {
  for_each = aws_subnet.private_subnets

  subnet_id       = each.value.id
  vpc_endpoint_id = aws_vpc_endpoint.kms_endpoint.id
}

resource "aws_vpc_endpoint" "sqs_endpoint" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${local.aws_region}.sqs"
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.vpc_endpoints.id,
  ]

  private_dns_enabled = true
}

resource "aws_vpc_endpoint_subnet_association" "private_sqs_association" {
  for_each = aws_subnet.private_subnets

  subnet_id       = each.value.id
  vpc_endpoint_id = aws_vpc_endpoint.sqs_endpoint.id
}
