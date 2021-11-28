variable "security_groups" {
  type = map(any)
  default = {
    "http" = 80,
    "https" = 443,
    "postgres" = 5432,
  }
}

resource "aws_security_group" "security_groups" {
  for_each = var.security_groups

  name        = format("allow_%s", each.key)
  description = format("Allow %s inbound traffic", upper(each.key))
  vpc_id      = aws_vpc.main.id

  ingress {
    description      = format("%s from VPC", upper(each.key))
    from_port        = each.value
    to_port          = each.value
    protocol         = "tcp"
    cidr_blocks      = [aws_vpc.main.cidr_block]
    ipv6_cidr_blocks = [aws_vpc.main.ipv6_cidr_block]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = format("allow_%s", each.key)
  }
}
