locals {
  repository_url = "ghcr.io/prixfixeco/api_server"
}

resource "aws_ecr_repository" "api_server" {
  name = "api_server"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
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
      "portMappings" : [
        {
          "containerPort" : 8888,
          "protocol" : "tcp",
        },
      ],
      "logConfiguration" : {
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
    container_port   = 8888
  }

  network_configuration {
    assign_public_ip = true

    security_groups = [
      aws_security_group.load_balancer.id,
    ]

    subnets = [for x in aws_subnet.private_subnets : x.id]
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

  managed_policy_arns = [
    aws_iam_policy.api_service_policy.arn,
  ]
}

resource "aws_iam_policy" "api_service_policy" {
  name = "api_service_policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sqs:SendMessage",
          "sqs:SendMessageBatch",
          "sqs:ReceiveMessage",
        ]
        Effect   = "Allow"
        Resource = aws_sqs_queue.writes_queue.arn
      },
    ]
  })
}