locals {
  log_services = ["base", "crm", "frontend"]
}

resource "aws_cloudwatch_log_group" "services" {
  for_each          = toset(local.log_services)
  name              = "/ecs/${var.project_name}/${each.key}"
  retention_in_days = 7

  tags = { Name = "/ecs/${var.project_name}/${each.key}" }
}
