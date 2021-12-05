# resource "aws_sqs_queue" "archives_queue" {
#   name       = "archives.fifo"
#   fifo_queue = true
# }

# resource "aws_ssm_parameter" "archives_queue_parameter" {
#   name  = "PRIXFIXE_ARCHIVES_QUEUE_URL"
#   type  = "String"
#   value = aws_sqs_queue.archives_queue.url
# }

# resource "aws_lambda_function" "archives_worker_lambda" {
#   function_name = "archives_worker"
#   handler       = "archives_worker"
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

# resource "aws_lambda_event_source_mapping" "archives_mapping" {
#   event_source_arn = aws_sqs_queue.archives_queue.arn
#   function_name    = aws_lambda_function.archives_worker_lambda.arn
# }