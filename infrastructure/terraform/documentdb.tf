# ─── DocumentDB (MongoDB-compatible) ─────────────────────────────────────────

resource "aws_security_group" "docdb" {
  name        = "${local.name}-docdb-sg"
  description = "DocumentDB access from ECS tasks"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 27017
    to_port         = 27017
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_http.id, aws_security_group.ecs_grpc.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-docdb-sg" }
}

resource "aws_docdb_subnet_group" "main" {
  name       = "${local.name}-docdb-subnet-group"
  subnet_ids = aws_subnet.database[*].id

  tags = { Name = "${local.name}-docdb-subnet-group" }
}

resource "aws_docdb_cluster_parameter_group" "main" {
  family      = "docdb5.0"
  name        = "${local.name}-docdb-params"
  description = "DocumentDB cluster parameter group"

  parameter {
    name  = "tls"
    value = "disabled"
  }

  tags = { Name = "${local.name}-docdb-params" }
}

resource "aws_docdb_cluster" "main" {
  cluster_identifier              = "${local.name}-docdb"
  engine                          = "docdb"
  master_username                 = "usegroadmin"
  master_password                 = var.mongodb_password
  backup_retention_period         = 1
  preferred_backup_window         = "02:00-03:00"
  skip_final_snapshot             = true
  db_subnet_group_name            = aws_docdb_subnet_group.main.name
  vpc_security_group_ids          = [aws_security_group.docdb.id]
  storage_encrypted               = true
  db_cluster_parameter_group_name = aws_docdb_cluster_parameter_group.main.name

  tags = { Name = "${local.name}-docdb" }
}

resource "aws_docdb_cluster_instance" "main" {
  identifier         = "${local.name}-docdb-0"
  cluster_identifier = aws_docdb_cluster.main.id
  instance_class     = "db.t3.medium"

  tags = { Name = "${local.name}-docdb-0" }
}
