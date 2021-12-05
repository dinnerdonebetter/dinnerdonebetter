# resource "aws_sqs_queue" "data_changes_queue" {
#   name       = "data_changes.fifo"
#   fifo_queue = true
# }

# resource "aws_ssm_parameter" "data_changes_queue_parameter" {
#   name  = "PRIXFIXE_DATA_CHANGES_QUEUE_URL"
#   type  = "String"
#   value = aws_sqs_queue.data_changes_queue.url
# }

# resource "aws_lambda_function" "data_changes_worker_lambda" {
#   function_name = "data_changes_worker"
#   handler       = "data_changes_worker"
#   role          = aws_iam_role.worker_lambda_role.arn
#   runtime       = local.lambda_runtime
#   memory_size   = local.memory_size
#   timeout       = local.timeout

#   tracing_config {
#     mode = "Active"
#   }

#   filename         = "writer_lambda.zip"
#   source_code_hash = filebase64sha256("writer_lambda.zip")
# }

# resource "aws_lambda_event_source_mapping" "data_changes_mapping" {
#   event_source_arn = aws_sqs_queue.data_changes_queue.arn
#   function_name    = aws_lambda_function.data_changes_worker_lambda.arn
# }