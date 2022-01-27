resource "aws_ecr_repository" "meal_plan_finalizer" {
  name = "meal_plan_finalizer"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_cloudwatch_log_group" "meal_plan_finalizer" {
  name              = "/ecs/meal_plan_finalizer"
  retention_in_days = local.log_retention_period_in_days
}

resource "aws_ecs_task_definition" "meal_plan_finalizer" {
  family = "meal_plan_finalizer"

  container_definitions = jsonencode([
    {
      name : "otel-collector",
      image : format("%s:latest", aws_ecr_repository.otel_collector.repository_url)
      essential : true,
      logConfiguration : {
        logDriver : "awslogs",
        options : {
          awslogs-region : local.aws_region,
          awslogs-group : "sidecars",
          awslogs-create-group : "true",
          awslogs-stream-prefix : "otel-collector"
        }
      }
    },
    {
      name  = "meal_plan_finalizer",
      image = format("%s:latest", aws_ecr_repository.meal_plan_finalizer.repository_url),
      essential : true,
      logConfiguration : {
        logDriver : "awslogs",
        options : {
          awslogs-region : local.aws_region,
          awslogs-group : aws_cloudwatch_log_group.meal_plan_finalizer.name,
          awslogs-stream-prefix : "ecs",
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

resource "aws_ecs_service" "meal_plan_finalizer" {
  name                               = "meal_plan_finalizer"
  task_definition                    = aws_ecs_task_definition.meal_plan_finalizer.arn
  cluster                            = aws_ecs_cluster.dev.id
  launch_type                        = "FARGATE"
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100
  desired_count                      = 1

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  network_configuration {
    assign_public_ip = true

    subnets = concat(
      [for x in aws_subnet.private_subnets : x.id],
    )
  }
}
