locals {
  api_database_username = "api_db_user"
}

resource "random_password" "api_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_sql_user" "api_user" {
  name     = local.api_database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.api_user_database_password.result
}
