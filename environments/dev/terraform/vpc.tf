resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/16"
  instance_tenancy                 = "default"
  assign_generated_ipv6_cidr_block = true
}

variable "availability_zones" {
    type = list(string)
    default = {
        "us-east-1a" = "10.0.1.0/26",
        "us-east-1b" = "10.0.1.1/26",
    }
}

resource "aws_subnet" "public_subnets" {
  for_each = toset(var.availability_zones)

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value
  availability_zone = each.key

  tags = {
    Name = format("main-%s", each.key)
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table_association" "PublicRTassociation" {
  subnet_id      = aws_subnet.east1a.id
  route_table_id = aws_route_table.public.id
}
