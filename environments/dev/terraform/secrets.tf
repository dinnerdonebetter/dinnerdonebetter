resource "random_string" "cookie_hash_key" {
  length  = 64
  special = false
}

resource "aws_ssm_parameter" "cookie_hash_key" {
  name  = "PRIXFIXE_COOKIE_HASH_KEY"
  type  = "SecureString"
  value = random_string.cookie_hash_key.result

  tags = merge(var.default_tags, {})
}

resource "random_string" "cookie_block_key" {
  length  = 64
  special = false
}

resource "aws_ssm_parameter" "cookie_block_key" {
  name  = "PRIXFIXE_COOKIE_BLOCK_KEY"
  type  = "SecureString"
  value = random_string.cookie_block_key.result

  tags = merge(var.default_tags, {})
}
