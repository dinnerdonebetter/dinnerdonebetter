resource "aws_ecr_repository" "dev-api-server" {
  name                 = "dev-api-server"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_lb_target_group" "dev-api-server" {
  name     = "dev-load-balancer"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_ecs_cluster" "dev-api-server" {
  name               = "dev-ecs-cluster"
  capacity_providers = ["FARGATE"]

  # configuration {
  #     log_configuration {
  #       cloud_watch_log_group_name     = aws_cloudwatch_log_group.example.name
  #     }
  # }

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# resource "aws_ecs_task_definition" "dev-api-server" {
#   family = "service"
#   container_definitions = jsonencode([
#     {
#       name      = "first"
#       image     = "dev_api_server"
#       cpu       = 10
#       memory    = 512
#       essential = true
#       portMappings = [
#         {
#           containerPort = 8080
#           hostPort      = 80
#         }
#       ]
#     }
#   ])
# }

# resource "aws_ecs_service" "dev-api-server" {
#   name            = "dev-api-server"
#   cluster         = aws_ecs_cluster.dev-api-server.id
#   task_definition = aws_ecs_task_definition.dev-api-server.arn
#   desired_count   = 3
#   iam_role        = data.aws_iam_role.worker_lambda_role.arn
#   depends_on      = [aws_iam_role_policy.foo]

#   ordered_placement_strategy {
#     type  = "binpack"
#     field = "cpu"
#   }

#   load_balancer {
#     target_group_arn = aws_lb_target_group.dev-api-server.arn
#     container_name   = "dev_api_server"
#     container_port   = 8080
#   }
# }