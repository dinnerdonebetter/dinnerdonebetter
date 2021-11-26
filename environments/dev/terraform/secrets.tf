resource "random_password" "cookie_hash_key" {
  length  = 64
  special = true
}

resource "aws_ssm_parameter" "cookie_hash_key" {
  name  = "PRIXFIXE_COOKIE_HASH_KEY"
  type  = "SecureString"
  value = random_password.cookie_hash_key.result
}

resource "random_password" "cookie_block_key" {
  length  = 64
  special = true
}

resource "aws_ssm_parameter" "cookie_block_key" {
  name  = "PRIXFIXE_COOKIE_BLOCK_KEY"
  type  = "SecureString"
  value = random_password.cookie_block_key.result
}
