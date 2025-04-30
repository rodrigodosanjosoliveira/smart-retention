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

output "cloudfront_cache_policy_no_cache_api_id" {
  description = "ID da política de cache desativado usada para rotas /api/*"
  value       = aws_cloudfront_cache_policy.no_cache_api.id
}

output "cloudfront_origin_request_policy_api_id" {
  description = "ID da política de request forwarding usada para rotas /api/*"
  value       = aws_cloudfront_origin_request_policy.api_origin_request.id
}
