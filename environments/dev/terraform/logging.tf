resource "aws_cloudwatch_log_group" "dev" {
  name              = "development"
  retention_in_days = 7
}
