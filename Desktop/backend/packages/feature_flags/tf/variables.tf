variable "image_version" {
  type        = string
  description = "The current image version"
}

variable "apigateway_id" {
  description = "API gateway ID"
  type        = string
}

variable "apigateway_stage_name" {
  description = "Stage name for HTTP apigateway configuration"
  type        = string
}

variable "aws_apigatewayv2_loop_gateway_id" {
  type        = string
  description = "API Gateway ID for the API that supports sonar loop"
}

variable "aws_apigatewayv2_loop_gateway_stage_name" {
  type        = string
  description = "API Gateway ID for the API that supports sonar loop"
}

variable "environment" {
  type        = string
  description = "Current environment of the project"
}

variable "lambda_image_uri" {
  type        = string
  description = "Image URI of the forms ECR repository"
}

variable "domain_name" {
  description = "Environments domain name"
  type        = string
}


variable "vpc_id" {
  type        = string
  description = "VPC to run select Lambdas inside"
}

variable "private_subnets" {
  type        = list(string)
  description = "List of private subnet IDs"
}

variable "rds_security_group" {
  type        = string
  description = "RDS security group to be used to give VPC access to RDS"
}

variable "external_security_group" {
  type        = string
  description = "S3 security group that allows the VPC to access external s3"
}

variable "db_username" {
  description = "Database administrator username"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "RDS root user password"
  type        = string
  sensitive   = true
}

variable "db_host" {
  description = "DB HOST"
  type        = string
}

variable "db_read_only_host" {
  description = "DB Read Only HOST"
  type        = string
}

variable "db_name" {
  description = "Name of form DB"
  type        = string
}

variable "db_port" {
  description = "DB port"
  type        = number
}