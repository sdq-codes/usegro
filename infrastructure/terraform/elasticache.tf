resource "aws_elasticache_subnet_group" "main" {
  name       = "${local.name}-redis-subnet-group"
  subnet_ids = aws_subnet.database[*].id

  tags = { Name = "${local.name}-redis-subnet-group" }
}

resource "aws_elasticache_replication_group" "main" {
  replication_group_id = "${local.name}-redis"
  description          = "Redis for ${local.name}"

  node_type            = var.redis_node_type
  num_cache_clusters   = 1
  parameter_group_name = "default.redis7"
  subnet_group_name    = aws_elasticache_subnet_group.main.name
  security_group_ids   = [aws_security_group.redis.id]

  engine_version             = "7.1"
  port                       = 6379
  at_rest_encryption_enabled = true
  transit_encryption_enabled  = true
  transit_encryption_mode     = "preferred"
  apply_immediately           = true

  automatic_failover_enabled = false
  multi_az_enabled           = false

  maintenance_window       = "sun:05:00-sun:06:00"
  snapshot_retention_limit = 1
  snapshot_window          = "04:00-05:00"

  tags = { Name = "${local.name}-redis" }
}
