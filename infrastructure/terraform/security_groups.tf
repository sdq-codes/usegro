# ─── ALB ────────────────────────────────────────────────────────────────────

resource "aws_security_group" "alb" {
  name        = "${local.name}-alb-sg"
  description = "Allow HTTP/HTTPS from the internet"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-alb-sg" }
}

# ─── ECS Tasks (HTTP services) ──────────────────────────────────────────────

resource "aws_security_group" "ecs_http" {
  name        = "${local.name}-ecs-http-sg"
  description = "ECS tasks reachable from ALB on port 8090 and 3000"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "HTTP from ALB (Go services)"
    from_port       = 8090
    to_port         = 8090
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  ingress {
    description     = "HTTP from ALB (whatsapp-gateway)"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-ecs-http-sg" }
}

# ─── ECS Tasks (gRPC internal services) ─────────────────────────────────────

resource "aws_security_group" "ecs_grpc" {
  name        = "${local.name}-ecs-grpc-sg"
  description = "Internal gRPC communication between ECS tasks"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "CRM gRPC"
    from_port   = 50051
    to_port     = 50051
    protocol    = "tcp"
    self        = true
  }

  ingress {
    description = "Messaging gRPC"
    from_port   = 50052
    to_port     = 50052
    protocol    = "tcp"
    self        = true
  }

  ingress {
    description = "Billing gRPC"
    from_port   = 50053
    to_port     = 50053
    protocol    = "tcp"
    self        = true
  }

  ingress {
    description = "Analytics gRPC"
    from_port   = 50054
    to_port     = 50054
    protocol    = "tcp"
    self        = true
  }

  # Also allow from HTTP services (crm calls gRPC services)
  ingress {
    description     = "gRPC from HTTP services"
    from_port       = 50051
    to_port         = 50054
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_http.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-ecs-grpc-sg" }
}

# ─── RDS ────────────────────────────────────────────────────────────────────

resource "aws_security_group" "rds" {
  name        = "${local.name}-rds-sg"
  description = "PostgreSQL access from ECS tasks"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_http.id, aws_security_group.ecs_grpc.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-rds-sg" }
}

# ─── ElastiCache ────────────────────────────────────────────────────────────

resource "aws_security_group" "redis" {
  name        = "${local.name}-redis-sg"
  description = "Redis access from ECS tasks"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_http.id, aws_security_group.ecs_grpc.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-redis-sg" }
}

