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