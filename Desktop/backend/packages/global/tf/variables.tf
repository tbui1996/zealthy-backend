variable "webapp_websocket_api_id" {
  type        = string
  description = "Websocket API Internal ID"
}

variable "webapp_websocket_api_execution_arn" {
  type        = string
  description = "Websocket API execution Internal ARN"
}

variable "webapp_http_api_id" {
  type        = string
  description = "HTTP API Internal ID"
}

variable "webapp_http_api_execution_arn" {
  type        = string
  description = "HTTP API execution Internal ARN"
}

variable "loop_websocket_api_id" {
  type        = string
  description = "Websocket API Internal ID"
}

variable "loop_websocket_api_execution_arn" {
  type        = string
  description = "Websocket API execution Internal ARN"
}

variable "loop_unconfirmed_websocket_api_id" {
  type = string
}

variable "loop_unconfirmed_websocket_api_execution_arn" {
  type = string
}

variable "loop_http_api_id" {
  type        = string
  description = "HTTP API Internal ID"
}

variable "loop_http_api_execution_arn" {
  type        = string
  description = "HTTP API execution Internal ARN"
}

variable "lambda_image_uri" {
  description = "Image URI of Docker lambda image for functions"
  type        = string
}

variable "image_version" {
  type        = string
  description = "The current image version"
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

variable "aws_cognito_user_pool_internals_jwks" {
  type        = string
  description = "Location of the JSON Web Key Set (JWKS) to be used to verify token signatures from internal cognito user pool"
}

variable "aws_cognito_user_pool_externals_jwks" {
  type        = string
  description = "Location of the JSON Web Key Set (JWKS) to be used to verify token signatures from external cognito user pool"
}

variable "name" {
  description = "Name of service"
  type        = string
}

variable "environment" {
  description = "Environment of the resource"
  type        = string
}

variable "group_names" {
  description = "Cognito group names"
  type        = set(string)
}

variable "group_routes" {
  description = "API Gateway routes that can be invoked for each group"
  type        = map(set(string))
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "vpc_public_subnets" {
  description = "VPC Public Subnets"
  type        = list(string)
}

variable "vpc_private_subnets" {
  description = "VPC Private Subnets"
  type        = list(string)
}

variable "domain_name" {
  description = "Environments domain name"
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

variable "db_name" {
  description = "Name of form DB"
  type        = string
}

variable "db_port" {
  description = "DB port"
  type        = number
}

variable "aws_account_id" {
  description = "AWS account ID"
  type        = string
}

variable "ecs_cluster_name" {
  type        = string
  description = "Name of ECS cluster"
}

variable "hosted_zone_id" {
  type        = string
  description = "Hosted zone ID for developer account"
}

variable "event_bus_arn" {
  type        = string
  description = "arn for EventBridge event bus"
}
