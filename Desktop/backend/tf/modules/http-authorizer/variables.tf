variable "authorizer_name" {
  type        = string
  description = "The name to be given to the authorizer"
}

variable "apigateway_id" {
  type        = string
  description = "The id of the api gateway to associate with authorizer"
}

variable "credentials_role_arn" {
  type        = string
  description = "Credentials arn of api gateway"
}

variable "lambda_invoke_arn" {
  type        = string
  description = "Invoke ARN of target lambda for authorizer"
}

variable "identity_sources" {
  type        = list(string)
  description = "Source of where authorization header will come from"
}

variable "statement_id" {
  type        = string
  description = "Lambda permission statement ID"
}

variable "lambda_function_name" {
  type        = string
  description = "Function name of the lambda to be invoked"
}

variable "api_gateway_execution_arn" {
  type        = string
  description = "Api gateway execution arn"
}