variable "route_key" {
  type        = string
  description = "The route key for this lambda integration"
}

variable "lambda_function" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "source_arn" {
  type = string
}

variable "api_id" {
  type = string
}

variable "authorizer_id" {
  description = "The authorizer id to be used when the authorization is CUSTOM or COGNITO_USER_POOLS"
  type        = string
  default     = null
}

variable "requires_auth" {
  description = "Boolean to indicate if authentication is required or not"
  type        = bool
  default     = false
}

variable "credentials_arn" {
  description = "Use to override aws_apigatewayv2_integration credential_arn"
  type        = string
  default     = ""
}
