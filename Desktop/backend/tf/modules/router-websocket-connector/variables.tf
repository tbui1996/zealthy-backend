variable "environment" {
  type        = string
  description = "Current environment of the application"
}

variable "name" {
  type        = string
  description = "Name of the package being worked in"
}

variable "lambda_image_uri" {
  type        = string
  description = "Image URI of the lambda"
}

variable "image_version" {
  type        = string
  description = "Image version represented by COMMIT_SHA"
}

variable "domain_name" {
  type        = string
  description = "Domain for route53 resources"
}

variable "loop_websocket_api_id" {
  type        = string
  description = "Websocket API ID"
}

variable "loop_websocket_api_execution_arn" {
  type        = string
  description = "Websocket API execution arn"
}

variable "router_dynamo_arns" {
  description = "Dynamo arns for tables and keys"
  type = object({
    websocket_connections         = string
    websocket_connections_kms_key = string
    pending_messages              = string
    pending_messages_kms_key      = string
  })
}

variable "lambda_receive" {
  type = object({
    function_name = string
  })
}