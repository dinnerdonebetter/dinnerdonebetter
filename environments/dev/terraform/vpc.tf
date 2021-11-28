resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/24"
  instance_tenancy                 = "default"
  assign_generated_ipv6_cidr_block = true

  tags = merge(var.default_tags, {})
}

resource "aws_subnet" "main" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.0.0/24"

  tags = merge(var.default_tags, {})
}