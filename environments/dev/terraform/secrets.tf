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
  length  = 64
  special = false
}

resource "aws_ssm_parameter" "cookie_block_key" {
  name   = "PRIXFIXE_COOKIE_BLOCK_KEY"
  type   = "SecureString"
  value  = random_string.cookie_block_key.result
  key_id = aws_kms_key.parameter_store_key.arn
}
