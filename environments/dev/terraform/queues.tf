resource "aws_sqs_queue" "writes_queue" {
  name       = "writes.fifo"
  fifo_queue = true

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_ssm_parameter" "writes_queue_parameter" {
  name  = "PRIXFIXE_WRITES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.writes_queue.url

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_event_source_mapping" "writes_mapping" {
  event_source_arn = aws_sqs_queue.writes_queue.arn
  function_name    = aws_lambda_function.writes_worker_lambda.arn
}

resource "aws_sqs_queue" "updates_queue" {
  name       = "updates.fifo"
  fifo_queue = true

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_ssm_parameter" "updates_queue_parameter" {
  name  = "PRIXFIXE_UPDATES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.updates_queue.url

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_event_source_mapping" "updates_mapping" {
  event_source_arn = aws_sqs_queue.updates_queue.arn
  function_name    = aws_lambda_function.updates_worker_lambda.arn
}

resource "aws_sqs_queue" "archives_queue" {
  name       = "archives.fifo"
  fifo_queue = true

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_ssm_parameter" "archives_queue_parameter" {
  name  = "PRIXFIXE_ARCHIVES_QUEUE_URL"
  type  = "String"
  value = aws_sqs_queue.archives_queue.url

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_event_source_mapping" "archives_mapping" {
  event_source_arn = aws_sqs_queue.archives_queue.arn
  function_name    = aws_lambda_function.archives_worker_lambda.arn
}

resource "aws_sns_topic" "data_changes_queue" {
  name = "data_changes"

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_ssm_parameter" "data_changes_queue_parameter" {
  name  = "PRIXFIXE_DATA_CHANGES_TOPIC_ID"
  type  = "String"
  value = aws_sns_topic.data_changes_queue.arn

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_lambda_event_source_mapping" "data_changes_mapping" {
  event_source_arn = aws_sns_topic.data_changes_queue.arn
  function_name    = aws_lambda_function.data_changes_worker_lambda.arn
}
