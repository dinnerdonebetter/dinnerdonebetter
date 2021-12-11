resource "aws_ecr_repository" "api_server" {
  name = "api_server"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}



resource "aws_security_group" "api_service" {
  name        = "prixfixe_api"
  description = "HTTP traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8000
    to_port     = 8000
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

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

resource "aws_security_group" "load_balancer" {
  name        = "load_balancer"
  description = "public internet traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8000
    to_port     = 8000
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

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

resource "aws_cloudwatch_log_group" "api_server" {
  name = "/ecs/api_server"
}

resource "aws_ecs_task_definition" "api_server" {
  family = "api_server"

  container_definitions = jsonencode([
    {
      name  = "api_server",
      image = format("%s:latest", aws_ecr_repository.api_server.repository_url),
      portMappings : [
        {
          "containerPort" : 8000,
          "hostPort" : 8000,
          "protocol" : "tcp",
        },
      ],
      healthCheck : {
        command : ["CMD-SHELL", "curl -f http://httpbin.org/get || exit 1"]
        interval : 5,
        retries : 4,
        startPeriod : 10,
      },
      logConfiguration : {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-region" : "us-east-1",
          "awslogs-group" : "/ecs/api_server",
          "awslogs-stream-prefix" : "ecs",
        },
      },
    },
  ])

  execution_role_arn = aws_iam_role.api_task_execution_role.arn
  task_role_arn      = aws_iam_role.api_task_role.arn

  # These are the minimum values for Fargate containers.
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]

  network_mode = "awsvpc"
}

resource "aws_ecs_cluster" "api" {
  name = "api_servers"
}

resource "aws_ecs_service" "api_server" {
  name            = "api_server"
  task_definition = aws_ecs_task_definition.api_server.arn
  cluster         = aws_ecs_cluster.api.id
  launch_type     = "FARGATE"

  desired_count = 1

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api_server"
    container_port   = 8000
  }

  network_configuration {
    assign_public_ip = true

    security_groups = [
      aws_security_group.api_service.id,
    ]

    subnets = concat(
      [for x in aws_subnet.public_subnets : x.id],
      [for x in aws_subnet.private_subnets : x.id],
    )
  }

  depends_on = [
    aws_alb_listener.api_http,
  ]
}

data "aws_iam_policy_document" "ecs_task_execution_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "api_task_execution_role" {
  name               = "api-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_execution_assume_role.json
}

# Normally we'd prefer not to hardcode an ARN in our Terraform, but since this is an AWS-managed policy, it's okay.
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.api_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "ecs_task_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "api_task_role" {
  name = "api-task-role"

  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json

  inline_policy {
    name   = "allow_sqs_queue_access"
    policy = data.aws_iam_policy_document.allow_to_manipulate_queues.json
  }

  inline_policy {
    name   = "allow_ssm_access"
    policy = data.aws_iam_policy_document.allow_parameter_store_access.json
  }

  inline_policy {
    name   = "allow_decrypt_ssm_parameters"
    policy = data.aws_iam_policy_document.allow_to_decrypt_parameters.json
  }
}

resource "cloudflare_record" "api_dot_prixfixe_dot_dev" {
  zone_id = var.CLOUDFLARE_ZONE_ID
  name    = "api"
  value   = aws_alb.api.dns_name
  type    = "CNAME"
  ttl     = 3600
}

output "alb_url" {
  value = "http://${aws_alb.api.dns_name}"
}