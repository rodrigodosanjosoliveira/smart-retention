output "backend_alb_dns" {
  description = "DNS público do Load Balancer (ALB) do backend Go"
  value       = aws_lb.backend.dns_name
}

output "cloudfront_distribution_url" {
  description = "URL pública do Frontend React via CloudFront"
  value       = "https://${aws_cloudfront_distribution.frontend.domain_name}"
}

output "s3_bucket_frontend" {
  description = "Nome do bucket S3 usado para o Frontend React"
  value       = aws_s3_bucket.frontend.bucket
}

output "rds_endpoint" {
  description = "Endpoint de conexão do banco RDS PostgreSQL"
  value       = aws_db_instance.postgresql.address
}

output "rds_port" {
  description = "Porta de conexão do banco RDS PostgreSQL"
  value       = aws_db_instance.postgresql.port
}
