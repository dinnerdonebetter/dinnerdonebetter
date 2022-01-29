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

