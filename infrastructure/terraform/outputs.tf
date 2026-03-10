output "alb_dns_name" {
  description = "Public DNS name of the Application Load Balancer"
  value       = aws_lb.main.dns_name
}

output "alb_zone_id" {
  description = "Route53 hosted zone ID for the ALB (use for ALIAS records)"
  value       = aws_lb.main.zone_id
}

output "rds_endpoint" {
  description = "RDS PostgreSQL endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "redis_primary_endpoint" {
  description = "ElastiCache Redis primary endpoint"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
  sensitive   = true
}

output "ecr_urls" {
  description = "ECR repository URLs for all services"
  value       = { for k, v in aws_ecr_repository.services : k => v.repository_url }
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "service_discovery_namespace" {
  description = "Internal DNS namespace for service discovery"
  value       = aws_service_discovery_private_dns_namespace.internal.name
}

output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}
