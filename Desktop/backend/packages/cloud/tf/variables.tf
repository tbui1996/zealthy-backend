variable "lambda_image_uri" {
  type        = string
  description = "ECR repo URI to upload imager to"
}

variable "image_version" {
  type        = string
  description = "The commit hash to set as the image version"
}

variable "domain_name" {
  type        = string
  description = "Domain name"
}

variable "environment" {
  type        = string
  description = "Environment"
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
  description = "DB NAME"
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

variable "last_production_commit" {
  type        = string
  description = "last production commit for versioning lambdas"
}

variable "doppler_host" {
  description = "Host string to Doppler DB"
  type        = string
  sensitive   = true
}

variable "doppler_port" {
  description = "Port for Doppler DB"
  type        = number
}

variable "doppler_user" {
  description = "User for logging in to Doppler DB"
  type        = string
  sensitive   = true
}

variable "doppler_pw" {
  description = "PW for logging in to Doppler DB"
  type        = string
  sensitive   = true
}

variable "doppler_dbname" {
  description = "Name of Doppler DB"
  type        = string
  sensitive   = true
}

variable "live_env" {
  description = "True if the environment is live environment, false if it is a developer environment"
  type        = bool
}
