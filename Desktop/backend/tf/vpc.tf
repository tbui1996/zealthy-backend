module "vpc" {
  source = "./modules/vpc"

  create_vpc           = !local.live_env
  private_subnet_names = var.private_subnet_names
  public_subnet_names  = var.public_subnet_names
  environment          = var.environment
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id       = module.vpc.vpc_id
  service_name = "com.amazonaws.${data.aws_region.current.name}.s3"

  route_table_ids = module.vpc.private_route_table_ids
}

resource "aws_vpc_endpoint" "dynamodb" {
  vpc_id       = module.vpc.vpc_id
  service_name = "com.amazonaws.${data.aws_region.current.name}.dynamodb"

  route_table_ids = module.vpc.private_route_table_ids
}

resource "aws_vpc_endpoint" "sqs" {
  vpc_id       = module.vpc.vpc_id
  service_name = "com.amazonaws.${data.aws_region.current.name}.sqs"

  vpc_endpoint_type  = "Interface"
  security_group_ids = [aws_security_group.external_service_security_group.id]

  subnet_ids          = module.vpc.private_subnets
  private_dns_enabled = true
}
