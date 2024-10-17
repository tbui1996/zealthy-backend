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
  description = "Image URI of the users ECR repository"
}

variable "domain_name" {
  type        = string
  description = "Domain name of hosted zone for application"
}

variable "loop_websocket_api_id" {
  type = string
}

variable "loop_unconfirmed_websocket_api_id" {
  description = "API ID for the loop side, unconfirmed web socket"
  type        = string
}


variable "router_dynamo_arns" {
  description = "Dynamo arns for tables and keys"
  type = object({
    websocket_connections                     = string
    websocket_connections_kms_key             = string
    pending_messages                          = string
    pending_messages_kms_key                  = string
    websocket_connections_internal            = string
    websocket_connections_internal_kms_key    = string
    unconfirmed_websocket_connections         = string
    unconfirmed_websocket_connections_kms_key = string
  })
}

# okta vars start
variable "okta_org_name" {
  type        = string
  description = "Okta org_name, e.g. dev-123456"
}

variable "okta_base_url" {
  type        = string
  description = "Okta base_url, e.g. okta.com"
}

variable "okta_api_token" {
  type        = string
  description = "Okta API token"
}

variable "callback_urls" {
  type        = list(string)
  description = "List of allowed callback urls for identity providers"
  default     = ["http://localhost:8080/dashboard/app"]
}
variable "logout_urls" {
  type        = list(string)
  description = "List of allowed logout urls for identity providers"
  default     = ["http://localhost:8080"]
}

variable "okta_user_id" {
  # only used for dev environments
  type        = string
  description = "Okta developer user id to associate with the app, e.g. 00u1jfzusmgDvUc1H5d7"
  default     = ""
}

variable "okta_username" {
  # only used for dev environments
  type        = string
  description = "Okta developer username/email to associate with the app, e.g. you@circulohealth.com"
  default     = ""
}
# okta vars end

variable "live_env" {
  type        = bool
  description = "True if the environment is live environment, false if it is a developer environment"
}
variable "internal_group_names" {
  type = set(string)
}

variable "external_group_names" {
  type = set(string)
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
