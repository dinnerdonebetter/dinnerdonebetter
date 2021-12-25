resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/16"
  instance_tenancy                 = "default"
  assign_generated_ipv6_cidr_block = true
  enable_dns_support               = true
  enable_dns_hostnames             = true

  tags = {
    Name = "dev"
  }
}

variable "public_availability_zones" {
  type = map(any)
  default = {
    "us-east-1a" = "10.0.1.0/26",
    "us-east-1b" = "10.0.1.64/26",
  }
}

resource "aws_subnet" "public_subnets" {
  for_each = var.public_availability_zones

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value
  availability_zone = each.key

  tags = {
    Name = format("public-%s", each.key)
  }
}

variable "private_availability_zones" {
  type = map(any)
  default = {
    "us-east-1a" = "10.0.129.0/26",
    "us-east-1b" = "10.0.129.64/26",
  }
}

resource "aws_subnet" "private_subnets" {
  for_each = var.private_availability_zones

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value
  availability_zone = each.key

  tags = {
    Name = format("private-%s", each.key)
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "main_internet_gateway"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "dev-public"
  }
}

resource "aws_route_table_association" "public_subnets" {
  for_each = aws_subnet.public_subnets

  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "dev-private"
  }
}

resource "aws_route_table_association" "private_subnets" {
  for_each = aws_subnet.private_subnets

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private.id
}

resource "aws_main_route_table_association" "main" {
  vpc_id         = aws_vpc.main.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route" "public_igw" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
}
