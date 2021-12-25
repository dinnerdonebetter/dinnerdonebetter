resource "aws_s3_bucket" "avatars_bucket" {
  bucket = "avatars.prixfixe.dev"
  acl    = "public-read"
}