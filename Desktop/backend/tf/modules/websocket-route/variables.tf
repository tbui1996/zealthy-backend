variable "route_key" {
  description = "Route key for route"
  type        = string
}

variable "requires_auth" {
  description = "Should attack authorizer"
  type        = bool
  default     = false
}

variable "authorizer_id" {
  description = "API Gateway V2 authorizer id"
  default     = null
}

variable "lambda_function" {
  description = "Lambda function"
  type = object({
    function_name = string
    invoke_arn    = string
  })
}

variable "websocket_api_id" {
  type        = string
  description = "Websocket API ID"
}

variable "websocket_api_execution_arn" {
  type        = string
  description = "Websocket API execution arn"
}
