variable "lambda_image_uri" {
  description = "Image URI of Docker lambda image for functions"
  type        = string
}

variable "image_version" {
  type        = string
  description = "The current image version"
}

variable "environment" {
  type = string
}

variable "domain_name" {
  type = string
}

variable "hosted_zone_id" {
  type = string
}

variable "name" {
  description = "Name of service"
  type        = string
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

variable "db_port" {
  description = "DB port"
  type        = number
}

variable "db_name" {
  description = "DB port"
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

variable "loop_websocket_api_id" {
  type        = string
  description = "API ID for the loop side, web socket"
}

variable "event_bus_arn" {
  description = "event bus arn"
  type        = string
}
