variable "domain_prefix" {
  description = "The prefix to your domain"
  type        = string
}

variable "domain_name" {
  description = "The name of the root domain"
  type        = string
}

variable "product_domain" {
  description = "A tag that specifies the project the domain is for"
  type        = string
}

variable "environment" {
  description = "Current environment of the project"
  type        = string
}

variable "api_id" {
  description = "APIGatewayV2 API ID"
  type        = string
}

variable "stage_id" {
  description = "APIGatewayV2 Stage ID"
  type        = string
}

variable "set_identifier" {
  description = "Identifies name and environment of A record"
  type        = string
}

variable "hosted_zone_id" {
  description = "Hosted zone associated with root domain"
  type        = string
}

