output "vpc_id" {
  value = var.create_vpc ? module.sonar_vpc[0].vpc_id : data.aws_vpc.vpc[0].id
}

output "public_subnets" {
  value = var.create_vpc ? module.sonar_vpc[0].public_subnets : data.aws_subnet_ids.public[0].ids
}

output "private_subnets" {
  value = var.create_vpc ? module.sonar_vpc[0].private_subnets : data.aws_subnet_ids.private[0].ids
}

output "private_route_table_ids" {
  value = var.create_vpc ? module.sonar_vpc[0].private_route_table_ids : data.aws_route_tables.private[0].ids
}

output "vpc_cidr_block" {
  value = var.create_vpc ? module.sonar_vpc[0].vpc_cidr_block : data.aws_vpc.vpc[0].cidr_block
}
