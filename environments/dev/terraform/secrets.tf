variable "SEGMENT_API_TOKEN" {}
variable "SENDGRID_API_TOKEN" {}

resource "aws_kms_key" "parameter_store_key" {
  description = "to encrypt parameter store secrets"
}

resource "random_string" "cookie_hash_key" {
  length  = 64
  special = false
}

resource "aws_ssm_parameter" "cookie_hash_key" {
  name   = "PRIXFIXE_COOKIE_HASH_KEY"
  type   = "SecureString"
  value  = random_string.cookie_hash_key.result
  key_id = aws_kms_key.parameter_store_key.arn
}

resource "random_string" "cookie_block_key" {
  length  = 32
  special = false
}

resource "aws_ssm_parameter" "cookie_block_key" {
  name   = "PRIXFIXE_COOKIE_BLOCK_KEY"
  type   = "SecureString"
  value  = random_string.cookie_block_key.result
  key_id = aws_kms_key.parameter_store_key.arn
}

resource "random_string" "paseto_local_key" {
  length  = 32
  special = false
}

resource "aws_ssm_parameter" "paseto_local_key" {
  name   = "PRIXFIXE_PASETO_LOCAL_MODE_KEY"
  type   = "SecureString"
  value  = random_string.paseto_local_key.result
  key_id = aws_kms_key.parameter_store_key.arn
}