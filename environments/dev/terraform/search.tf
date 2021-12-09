resource "aws_cloudwatch_log_resource_policy" "dev" {
  policy_name = "dev"

  policy_document = <<CONFIG
 {
   "Version": "2012-10-17",
   "Statement": [
     {
       "Effect": "Allow",
       "Principal": {
         "Service": "es.amazonaws.com"
       },
       "Action": [
         "logs:PutLogEvents",
         "logs:PutLogEventsBatch",
         "logs:CreateLogStream"
       ],
       "Resource": "arn:aws:logs:*"
     }
   ]
 }
 CONFIG
}

resource "aws_iam_service_linked_role" "es" {
  aws_service_name = "es.amazonaws.com"
}

resource "aws_elasticsearch_domain" "search" {
  domain_name           = "dev-search"
  elasticsearch_version = "2.3"

  cluster_config {
    instance_type = "t2.micro.elasticsearch"
  }

  ebs_options {
    ebs_enabled = true
    volume_size = 10
  }

  log_publishing_options {
    cloudwatch_log_group_arn = aws_cloudwatch_log_group.dev.arn
    log_type                 = "INDEX_SLOW_LOGS"
  }

  vpc_options {
    subnet_ids = [aws_subnet.private_subnets["us-east-1a"].id]

    security_group_ids = [
      aws_security_group.search.id,
    ]
  }

  depends_on = [aws_iam_service_linked_role.es]
}

resource "aws_ssm_parameter" "search_url" {
  name  = "PRIXFIXE_ELASTICSEARCH_INSTANCE_URL"
  type  = "String"
  value = aws_elasticsearch_domain.search.arn
}