variable "aws_region" {
  default = "us-east-1"
}

variable "project_name" {
  description = "Prefixo para recursos AWS"
  default     = "smart-retention"
}

variable "db_username" {
  description = "Username para o banco de dados PostgreSQL"
}

variable "db_password" {
  description = "Senha para o banco de dados PostgreSQL"
  sensitive   = true
}