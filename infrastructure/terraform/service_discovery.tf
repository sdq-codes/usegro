# Private DNS namespace for internal service-to-service communication

resource "aws_service_discovery_private_dns_namespace" "internal" {
  name        = "${var.project_name}.internal"
  description = "Internal service discovery for ${local.name}"
  vpc         = aws_vpc.main.id

  tags = { Name = "${local.name}-namespace" }
}

resource "aws_service_discovery_service" "crm" {
  name = "crm"

  dns_config {
    namespace_id   = aws_service_discovery_private_dns_namespace.internal.id
    routing_policy = "MULTIVALUE"
    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}
