resource "aws_ecr_repository" "dev_api_server" {
  name                 = "dev_api_server"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = merge(var.default_tags, {})
}

resource "aws_lb_target_group" "dev_api_server" {
  name     = "load_balancer"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
}

resource "aws_ecs_cluster" "dev_api_server" {
  name               = "dev"
  capacity_providers = "FARGATE"

  configuration {
    logging = "DEFAULT"
  }

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = merge(var.default_tags, {})
}

# resource "aws_ecs_task_definition" "dev_api_server" {
#   family = "service"
#   container_definitions = jsonencode([
#     {
#       name      = "first"
#       image     = "debugserver"
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

# resource "aws_ecs_service" "dev_api_server" {
#   name            = "dev_api_server"
#   cluster         = aws_ecs_cluster.dev_api_server.id
#   task_definition = aws_ecs_task_definition.dev_api_server.arn
#   desired_count   = 3
#   iam_role        = data.aws_iam_role.worker_lambda_role.arn
#   depends_on      = [aws_iam_role_policy.foo]

#   ordered_placement_strategy {
#     type  = "binpack"
#     field = "cpu"
#   }

#   load_balancer {
#     target_group_arn = aws_lb_target_group.dev_api_server.arn
#     container_name   = "dev_api_server"
#     container_port   = 8080
#   }
# }