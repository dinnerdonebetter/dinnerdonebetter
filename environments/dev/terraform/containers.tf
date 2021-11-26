resource "aws_ecr_repository" "dev_api_server" {
  name                 = "dev_api_server"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = merge(var.default_tags, {})
}