locals {
  lambda_runtime = "go1.x"
  lambda_handler = "main"
  memory_size    = 128
  timeout        = 15

  cloudwatch_exclude_lambda_events_filter = "[timestamp != \"START\" && timestamp != \"END \" &&timestamp != \"REPORT\" && timestamp != \"RequestId: \", aws_request_id, log]"
}