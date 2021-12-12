locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 128
  timeout        = 15

  cloudwatch_exclude_lambda_events_filter = "-\"START RequestId\" -\"END RequestId\" -\"REPORT RequestId\""
}