locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 1024
  timeout        = 30
}

module "writes_worker_lambda" {
  source = "terraform-aws-modules/lambda/aws//modules/docker-build"

  create_ecr_repo      = true
  ecr_repo             = "writes_worker"
  image_tag            = "latest"
  image_tag_mutability = "MUTABLE"
  docker_file_path     = "./dockerfiles/writes_worker.Dockerfile"
}

resource "aws_lambda_function" "writes_worker_lambda" {
  function_name = "writes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = module.writes_worker_lambda.image_uri
}

module "updates_worker_lambda" {
  source = "terraform-aws-modules/lambda/aws//modules/docker-build"

  create_ecr_repo      = true
  ecr_repo             = "updates_worker"
  image_tag            = "latest"
  image_tag_mutability = "MUTABLE"
  docker_file_path     = "./dockerfiles/updates_worker.Dockerfile"
}

resource "aws_lambda_function" "updates_worker_lambda" {
  function_name = "updates_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = module.updates_worker_lambda.image_uri
}

module "archives_worker_lambda" {
  source = "terraform-aws-modules/lambda/aws//modules/docker-build"

  create_ecr_repo      = true
  ecr_repo             = "archives_worker"
  image_tag            = "latest"
  image_tag_mutability = "MUTABLE"
  docker_file_path     = "./dockerfiles/archives_worker.Dockerfile"
}

resource "aws_lambda_function" "archives_worker_lambda" {
  function_name = "archives_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = module.archives_worker_lambda.image_uri
}

module "data_changes_worker_lambda" {
  source = "terraform-aws-modules/lambda/aws//modules/docker-build"

  create_ecr_repo      = true
  ecr_repo             = "data_changes_worker"
  image_tag            = "latest"
  image_tag_mutability = "MUTABLE"
  docker_file_path     = "./dockerfiles/data_changes_worker.Dockerfile"
}

resource "aws_lambda_function" "data_changes_worker_lambda" {
  function_name = "data_changes_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = module.data_changes_worker_lambda.image_uri
}

module "chores_worker_lambda" {
  source = "terraform-aws-modules/lambda/aws//modules/docker-build"

  create_ecr_repo      = true
  ecr_repo             = "chores_worker"
  image_tag            = "latest"
  image_tag_mutability = "MUTABLE"
  docker_file_path     = "./dockerfiles/chores_worker.Dockerfile"
}

resource "aws_lambda_function" "chores_worker_lambda" {
  function_name = "chores_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  memory_size   = local.memory_size
  timeout       = local.timeout
  package_type  = "Image"

  tracing_config {
    mode = "Active"
  }

  image_uri = module.chores_worker_lambda.image_uri
}

resource "aws_cloudwatch_event_rule" "every_minute" {
  name                = "every-minute"
  description         = "Fires every minute"
  schedule_expression = "rate(1 minute)"
}

resource "aws_cloudwatch_event_target" "run_chores_every_minute" {
  rule      = aws_cloudwatch_event_rule.every_minute.name
  target_id = "chores_worker"
  arn       = aws_lambda_function.chores_worker_lambda.arn
}