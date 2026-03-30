locals {
  common_go_env = [
    { name = "APP_ENV", value = "production" },
    { name = "POSTGRES_HOST", value = aws_db_instance.main.address },
    { name = "POSTGRES_PORT", value = "5432" },
    { name = "POSTGRES_DB", value = var.db_name },
    { name = "POSTGRES_USER", value = var.db_username },
    { name = "POSTGRES_SSL_MODE", value = "require" },
    { name = "REDIS_HOST", value = aws_elasticache_replication_group.main.primary_endpoint_address },
    { name = "REDIS_PORT", value = "6379" },
    { name = "AWS_REGION", value = var.aws_region },
    { name = "SES_FROM_EMAIL", value = var.ses_from_email },
  ]

  common_go_secrets = [
    { name = "POSTGRES_PASSWORD", valueFrom = aws_secretsmanager_secret.app["db_password"].arn },
    { name = "AUTH_API_SECRET", valueFrom = aws_secretsmanager_secret.app["auth_api_secret"].arn },
    { name = "SENTRY_DSN", valueFrom = aws_secretsmanager_secret.app["sentry_dsn"].arn },
  ]

  base_extra_secrets = [
    { name = "GOOGLE_CLIENT_SECRET", valueFrom = aws_secretsmanager_secret.app["google_client_secret"].arn },
    { name = "FACEBOOK_APP_SECRET", valueFrom = aws_secretsmanager_secret.app["facebook_app_secret"].arn },
  ]
}

# ─── ECS Cluster ─────────────────────────────────────────────────────────────

resource "aws_ecs_cluster" "main" {
  name = "${local.name}-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled"
  }

  tags = { Name = "${local.name}-cluster" }
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name       = aws_ecs_cluster.main.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
  }
}

# ─── Base Service ─────────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "base" {
  family                   = "${local.name}-base"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name      = "base"
    image     = "${aws_ecr_repository.services["base"].repository_url}:production"
    essential = true

    portMappings = [{ containerPort = 8090, protocol = "tcp" }]

    environment = concat(local.common_go_env, [
      { name = "HTTP_PORT", value = "8090" },
      { name = "GOOGLE_CLIENT_ID", value = var.google_client_id },
      { name = "FACEBOOK_APP_ID", value = var.facebook_app_id },
      { name = "FRONTEND_URL", value = var.domain_name != "" ? "https://${var.domain_name}" : "" },
    ])

    secrets = concat(local.common_go_secrets, local.base_extra_secrets)

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.services["base"].name
        awslogs-region        = var.aws_region
        awslogs-stream-prefix = "ecs"
      }
    }
  }])

  tags = { Name = "${local.name}-base" }
}

resource "aws_ecs_service" "base" {
  name            = "${local.name}-base"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.base.arn
  desired_count   = var.base_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.ecs_http.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.base.arn
    container_name   = "base"
    container_port   = 8090
  }

  depends_on = [aws_lb_listener.http]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  tags = { Name = "${local.name}-base" }
}

# ─── CRM Service ─────────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "crm" {
  family                   = "${local.name}-crm"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name      = "crm"
    image     = "${aws_ecr_repository.services["crm"].repository_url}:production"
    essential = true

    portMappings = [
      { containerPort = 8090, protocol = "tcp" },
      { containerPort = 50051, protocol = "tcp" },
    ]

    environment = concat(local.common_go_env, [
      { name = "HTTP_PORT", value = "8090" },
      { name = "GRPC_PORT", value = "50051" },
      { name = "DYNAMO_ENDPOINT", value = "" },
      { name = "DYNAMO_REGION", value = var.aws_region },
      { name = "DYNAMO_FORM_TABLE_NAME", value = "forms" },
      { name = "FRONTEND_URL", value = var.domain_name != "" ? "https://${var.domain_name}" : "" },
      { name = "ALB_URL", value = "http://${aws_lb.main.dns_name}" },
    ])

    secrets = local.common_go_secrets

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.services["crm"].name
        awslogs-region        = var.aws_region
        awslogs-stream-prefix = "ecs"
      }
    }
  }])

  tags = { Name = "${local.name}-crm" }
}

resource "aws_ecs_service" "crm" {
  name            = "${local.name}-crm"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.crm.arn
  desired_count   = var.crm_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.ecs_http.id, aws_security_group.ecs_grpc.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.crm.arn
    container_name   = "crm"
    container_port   = 8090
  }

  service_registries {
    registry_arn = aws_service_discovery_service.crm.arn
  }

  depends_on = [aws_lb_listener.http]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  tags = { Name = "${local.name}-crm" }
}

# ─── Catalog Service ──────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "catalog" {
  family                   = "${local.name}-catalog"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name      = "catalog"
    image     = "${aws_ecr_repository.services["catalog"].repository_url}:production"
    essential = true

    portMappings = [
      { containerPort = 8090, protocol = "tcp" },
      { containerPort = 50051, protocol = "tcp" },
    ]

    environment = concat(local.common_go_env, [
      { name = "HTTP_PORT", value = "8090" },
      { name = "GRPC_PORT", value = "50051" },
      { name = "FRONTEND_URL", value = var.domain_name != "" ? "https://${var.domain_name}" : "" },
      { name = "ALB_URL", value = "http://${aws_lb.main.dns_name}" },
    ])

    secrets = local.common_go_secrets

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.services["catalog"].name
        awslogs-region        = var.aws_region
        awslogs-stream-prefix = "ecs"
      }
    }
  }])

  tags = { Name = "${local.name}-catalog" }
}

resource "aws_ecs_service" "catalog" {
  name            = "${local.name}-catalog"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.catalog.arn
  desired_count   = var.catalog_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.ecs_http.id, aws_security_group.ecs_grpc.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.catalog.arn
    container_name   = "catalog"
    container_port   = 8090
  }

  depends_on = [aws_lb_listener.http]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  tags = { Name = "${local.name}-catalog" }
}

# ─── Billing Service ──────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "billing" {
  family                   = "${local.name}-billing"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name      = "billing"
    image     = "${aws_ecr_repository.services["billing"].repository_url}:production"
    essential = true

    portMappings = [{ containerPort = 8090, protocol = "tcp" }]

    environment = concat(local.common_go_env, [
      { name = "HTTP_PORT", value = "8090" },
    ])

    secrets = local.common_go_secrets

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.services["billing"].name
        awslogs-region        = var.aws_region
        awslogs-stream-prefix = "ecs"
      }
    }
  }])

  tags = { Name = "${local.name}-billing" }
}

resource "aws_ecs_service" "billing" {
  name            = "${local.name}-billing"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.billing.arn
  desired_count   = var.billing_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.ecs_http.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.billing.arn
    container_name   = "billing"
    container_port   = 8090
  }

  depends_on = [aws_lb_listener.http]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  tags = { Name = "${local.name}-billing" }
}

# ─── Frontend Service ─────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "frontend" {
  family                   = "${local.name}-frontend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name      = "frontend"
    image     = "${aws_ecr_repository.services["frontend"].repository_url}:production"
    essential = true

    portMappings = [{ containerPort = 3000, protocol = "tcp" }]

    environment = [
      { name = "NODE_ENV", value = "production" },
      { name = "NUXT_PUBLIC_API_BASE", value = "http://${aws_lb.main.dns_name}/api/v1" },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = aws_cloudwatch_log_group.services["frontend"].name
        awslogs-region        = var.aws_region
        awslogs-stream-prefix = "ecs"
      }
    }
  }])

  tags = { Name = "${local.name}-frontend" }
}

resource "aws_ecs_service" "frontend" {
  name            = "${local.name}-frontend"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.frontend.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.ecs_http.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.frontend.arn
    container_name   = "frontend"
    container_port   = 3000
  }

  depends_on = [aws_lb_listener.http]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  tags = { Name = "${local.name}-frontend" }
}
