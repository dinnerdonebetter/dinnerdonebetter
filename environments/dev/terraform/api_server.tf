locals {
  repository_url = "ghcr.io/prixfixeco/api_server"
}

resource "aws_ecr_repository" "api_server" {
  name = "api_server"
}

resource "aws_cloudwatch_log_group" "api" {
  name = "/ecs/api"
}

resource "aws_ecs_task_definition" "api_server" {
  family = "api"

  container_definitions = <<EOF
  [
    {
      "name": "api",
      "image": "${aws_ecr_repository.api_server.repository_url}:latest",
      "portMappings": [
        {
          "containerPort": 80
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-region": "us-east-1",
          "awslogs-group": "/ecs/api",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
EOF

  execution_role_arn = aws_iam_role.api_task_execution_role.arn

  # These are the minimum values for Fargate containers.
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]

  network_mode = "awsvpc"
}

resource "aws_ecs_cluster" "api" {
  name = "api"
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
    container_port   = 80
  }

  network_configuration {
    assign_public_ip = true

    security_groups = [
      aws_security_group.egress_all.id,
      aws_security_group.allow_http.id,
      aws_security_group.allow_https.id,
    ]

    subnets = [for x in aws_subnet.private_subnets : x.id]
  }

  depends_on = [
    aws_alb_listener.api_http,
  ]
}

resource "aws_iam_role" "api_task_execution_role" {
  name               = "api-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json
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

# Normally we'd prefer not to hardcode an ARN in our Terraform, but since this is an AWS-managed policy, it's okay.
data "aws_iam_policy" "ecs_task_execution_role" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.api_task_execution_role.name
  policy_arn = data.aws_iam_policy.ecs_task_execution_role.arn
}
