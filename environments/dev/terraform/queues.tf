resource "aws_sqs_queue" "writes_queue" {
  name       = "writes.fifo"
  fifo_queue = true

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_sqs_queue" "updates_queue" {
  name       = "updates.fifo"
  fifo_queue = true

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_sqs_queue" "archives_queue" {
  name       = "archives.fifo"
  fifo_queue = true

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}

resource "aws_sns_topic" "data_changes_queue" {
  name = "data_changes"

  tags = {
    Environment = "dev"
    Terraform   = "true"
  }
}
