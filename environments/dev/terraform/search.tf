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
  value = format("http://%s", aws_elasticsearch_domain.search.endpoint)
}

resource "aws_security_group" "search" {
  name        = "elasticsearch"
  description = "Elasticsearch traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 9200
    to_port     = 9200
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/0"]
  }
}
