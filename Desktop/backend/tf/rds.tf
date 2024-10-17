resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "sonar_rds_subnet_group"
  subnet_ids = !local.live_env ? module.vpc.public_subnets : module.vpc.private_subnets

  tags = {
    Name        = "Sonar RDS"
    Environment = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "rds_security_group" {
  name   = "sonar_rds_security_group"
  vpc_id = module.vpc.vpc_id

  description = "Security group for sonar-service-rds-cluster to allow internal access to main sonar databased hosted on RDS"

  ingress {
    description = "Allow VPC to reach DB"
    from_port   = 5432
    to_port     = 5432

    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow DB to respond to VPC"
    from_port   = 5432
    to_port     = 5432

    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "Sonar RDS"
    Environment = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_kms_key" "sonar_service_rds_cluster_kms_key" {
  description             = "KMS key for sonar-service-rds-cluster RDS cluster"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_rds_cluster" "sonar_service_rds_cluster" {
  cluster_identifier      = "sonar-service-rds-cluster"
  engine                  = "aurora-postgresql"
  database_name           = "sonar"
  master_username         = "circulo"
  master_password         = var.db_password
  vpc_security_group_ids  = [aws_security_group.rds_security_group.id]
  db_subnet_group_name    = aws_db_subnet_group.rds_subnet_group.name
  backup_retention_period = 5
  skip_final_snapshot     = true
  preferred_backup_window = "07:00-09:00"
  kms_key_id              = aws_kms_key.sonar_service_rds_cluster_kms_key.arn
  storage_encrypted       = true
}

#tfsec:ignore:aws-rds-no-public-db-access
resource "aws_rds_cluster_instance" "sonar_service_rds_instance" {
  count               = 1
  identifier          = "sonar-service-rds-instance-${count.index}"
  cluster_identifier  = aws_rds_cluster.sonar_service_rds_cluster.id
  instance_class      = "db.t3.large"
  engine              = aws_rds_cluster.sonar_service_rds_cluster.engine
  engine_version      = aws_rds_cluster.sonar_service_rds_cluster.engine_version
  publicly_accessible = !local.live_env
  apply_immediately   = true
}

# for adding external project's migrations to our local dev environment.
provider "postgresql" {
  host            = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  port            = aws_rds_cluster.sonar_service_rds_cluster.port
  database        = "postgres"
  username        = aws_rds_cluster.sonar_service_rds_cluster.master_username
  password        = var.db_password
  sslmode         = "require"
  connect_timeout = 15
}

# for adding external project's migrations to our local dev environment.
resource "postgresql_database" "external" {
  count = local.live_env ? 0 : 1
  name  = "external"
  owner = aws_rds_cluster.sonar_service_rds_cluster.master_username
}
