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

variable "form_db_name" {
  description = "Name of form DB"
  type        = string
}

variable "db_port" {
  description = "DB port"
  type        = number
}

variable "image_version" {
  type        = string
  description = "The current image version"
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
  type        = string
  description = "Domain name of hosted zone for application"
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

variable "receive_queue_arn" {
  type = string
}

variable "receive_queue_kms_arn" {
  type = string
}

variable "send_queue_arn" {
  type = string
}

variable "send_queue_kms_arn" {
  type = string
}
