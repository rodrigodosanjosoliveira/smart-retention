resource "aws_ssm_parameter" "app_env" {
  name  = "/smart-retention/APP_ENV"
  type  = "String"
  value = "production"
}

resource "aws_ssm_parameter" "db_host" {
  name  = "/smart-retention/DB_HOST"
  type  = "String"
  value = aws_db_instance.postgresql.address
}

resource "aws_ssm_parameter" "db_user" {
  name  = "/smart-retention/DB_USER"
  type  = "String"
  value = "smartretention"
}

resource "aws_ssm_parameter" "db_password" {
  name  = "/smart-retention/DB_PASSWORD"
  type  = "SecureString"
  value = var.db_password
}

resource "aws_ssm_parameter" "db_name" {
  name  = "/smart-retention/DB_NAME"
  type  = "String"
  value = "smartretention"
}
