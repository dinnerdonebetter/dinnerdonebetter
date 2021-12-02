resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/16"
  instance_tenancy                 = "default"
  assign_generated_ipv6_cidr_block = true

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

# resource "aws_internet_gateway" "main" {
#   vpc_id = aws_vpc.main.id

#   tags = {
#     Name = "main_internet_gateway"
#   }
# }

# resource "aws_route_table" "public" {
#   vpc_id = aws_vpc.main.id

#   route {
#     cidr_block = "0.0.0.0/0"
#     gateway_id = aws_internet_gateway.main.id
#   }
# }

# resource "aws_route_table_association" "public_subnets" {
#   for_each = aws_subnet.public_subnets

#   subnet_id      = each.value.id
#   route_table_id = aws_route_table.public.id
# }

# resource "aws_route_table" "private" {
#   vpc_id = aws_vpc.main.id

#   route {
#     cidr_block = "0.0.0.0/0"
#   }
# }

# resource "aws_route_table_association" "private_subnets" {
#   for_each = aws_subnet.private_subnets

#   subnet_id      = each.value.id
#   route_table_id = aws_route_table.private.id
# }

# resource "aws_route" "public_igw" {
#   route_table_id         = aws_route_table.public.id
#   destination_cidr_block = "0.0.0.0/0"
#   gateway_id             = aws_internet_gateway.main.id
# }

resource "aws_alb" "api" {
  name               = "api-lb"
  internal           = false
  load_balancer_type = "application"

  subnets = aws_subnet.public_subnets.*.id

  security_groups = [
    aws_security_group.allow_http.id,
    aws_security_group.allow_https.id,
    aws_security_group.egress_all.id,
  ]

  # depends_on = [aws_internet_gateway.main]
}

resource "aws_lb_target_group" "api" {
  name        = "api"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id

  depends_on = [aws_alb.api]
}


resource "aws_alb_listener" "api_http" {
  load_balancer_arn = aws_alb.api.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

output "alb_url" {
  value = "http://${aws_alb.api.dns_name}"
}


# resource "aws_alb_listener" "api_https" {
#   load_balancer_arn = aws_alb.api.arn
#   port              = "443"
#   protocol          = "HTTPS"
#   ssl_policy        = "ELBSecurityPolicy-2016-08"
#   certificate_arn   = aws_acm_certificate.api.arn

#   default_action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.api.arn
#   }
# }

# resource "aws_acm_certificate" "api" {
#   domain_name       = "api.prixfixe.dev"
#   validation_method = "DNS"

#   options {
#     certificate_transparency_logging_preference = "ENABLED"
#   }

#   lifecycle {
#     create_before_destroy = true
#   }

#   tags = {
#     Name = "dev_api"
#   }
# }

# output "domain_validations" {
#   value = aws_acm_certificate.api.domain_validation_options
# }