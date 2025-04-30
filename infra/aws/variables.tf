variable "environment" {
  description = "Nome do ambiente (ex: production)"
  type        = string
}

variable "aws_region" {
  default = "us-east-1"
}

variable "project_name" {
  description = "Prefixo para recursos AWS"
  default     = "smart-retention"
}

variable "db_username" {
  description = "Username para o banco de dados PostgreSQL"
    type        = string
    default     = "postgres"
}

variable "db_password" {
  description = "Senha para o banco de dados PostgreSQL"
  type = string
  sensitive   = true
    default     = "Precious78" # Altere para uma senha segura
}

variable "db_name" {
  description = "Nome do banco de dados PostgreSQL"
  default     = "smartretention"
}