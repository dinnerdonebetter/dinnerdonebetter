resource "aws_sqs_queue" "meal_planning_prober_dead_letter" {
  name                    = "meal_planning_prober_dead_letter"
  sqs_managed_sse_enabled = true
}

resource "aws_sqs_queue" "meal_planning_prober_queue" {
  name                    = "meal_planning_prober"
  sqs_managed_sse_enabled = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.meal_planning_prober_dead_letter.arn
    maxReceiveCount     = 1
  })
}

resource "aws_lambda_function" "meal_planning_prober_worker_lambda" {
  function_name = "prober_meal_planning"
  handler       = "prober_meal_planning"
  role          = aws_iam_role.worker_lambda_role.arn
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = 300

  filename = data.archive_file.dummy_zip.output_path

  depends_on = [
    aws_cloudwatch_log_group.meal_planning_prober_worker_lambda_logs,
  ]
}

resource "aws_cloudwatch_event_rule" "every_five_minutes" {
  name                = "every-five-minutes"
  description         = "Fires every five minutes"
  schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "run_meal_planning_prober_every_five_minutes" {
  rule      = aws_cloudwatch_event_rule.every_five_minutes.name
  target_id = "meal_planning_prober_worker"
  arn       = aws_lambda_function.meal_planning_prober_worker_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_meal_planning_prober_worker" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.meal_planning_prober_worker_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_five_minutes.arn
}

resource "aws_cloudwatch_log_group" "meal_planning_prober_worker_lambda_logs" {
  name              = "/aws/lambda/meal_planning_prober_worker"
  retention_in_days = local.log_retention_period_in_days
}
