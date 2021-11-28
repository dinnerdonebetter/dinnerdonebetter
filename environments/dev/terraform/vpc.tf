resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/24"
  instance_tenancy                 = "default"
  assign_generated_ipv6_cidr_block = true

  tags = merge(var.default_tags, {})
}

resource "aws_subnet" "east1a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "us-east-1a"

  tags = merge(var.default_tags, {
    Name = "main"
  })
}

resource "aws_subnet" "east1b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "us-east-1b"

  tags = merge(var.default_tags, {
    Name = "main"
  })
}