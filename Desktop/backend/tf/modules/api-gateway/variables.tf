variable "gateway_name" {
  type        = string
  description = "Name of the gateway"
}

variable "domain_name" {
  type        = string
  description = "Domain name of route53 configuration"
}

variable "environment" {
  type        = string
  description = "Current environment"
}

variable "deployment_triggers" {
  type        = list(string)
  description = "Routes that redeploy api gateway"
}

variable "aws_region" {
  type        = string
  description = "AWS region"
}

variable "aws_account_id" {
  type        = string
  description = "AWS Account ID"
}

variable "acm_certificate_arn" {
  type        = string
  description = "Certificate ARN of route53 configuration"
}

variable "api_mapping_key" {
  type        = string
  description = "The route mapping key to associate with the new api gateway"
}