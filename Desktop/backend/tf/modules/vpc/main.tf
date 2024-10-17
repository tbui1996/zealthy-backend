terraform {
  required_version = ">= 0.15.0"
}

data "aws_region" "current" {}

module "sonar_vpc" {
  count  = var.create_vpc ? 1 : 0
  source = "terraform-aws-modules/vpc/aws"

  name = "sonar_vpc"
  cidr = "10.0.0.0/16"

  azs             = ["us-east-2a", "us-east-2b", "us-east-2c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway   = true
  enable_vpn_gateway   = true
  enable_dns_hostnames = true

  tags = {
    Environment = var.environment
  }
}

data "aws_ssm_parameter" "vpc_id" {
  count = var.create_vpc ? 0 : 1
  name  = "/vpc/id"
}

data "aws_vpc" "vpc" {
  count = var.create_vpc ? 0 : 1
  id    = data.aws_ssm_parameter.vpc_id[0].value
}

data "aws_subnet_ids" "private" {
  count = var.create_vpc ? 0 : 1

  vpc_id = data.aws_vpc.vpc[0].id

  filter {
    name   = "tag:Name"
    values = var.private_subnet_names
  }
}

data "aws_subnet_ids" "public" {
  count = var.create_vpc ? 0 : 1

  vpc_id = data.aws_vpc.vpc[0].id

  filter {
    name   = "tag:Name"
    values = var.public_subnet_names
  }
}

data "aws_route_tables" "private" {
  count = var.create_vpc ? 0 : 1

  vpc_id = data.aws_vpc.vpc[0].id

  filter {
    name = "tag:Name"
    // Subnet names are the same as route table names
    values = var.private_subnet_names
  }
}
