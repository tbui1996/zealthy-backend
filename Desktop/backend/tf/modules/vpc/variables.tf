variable "create_vpc" {
  description = "Boolean to create a VPC or use exisiting (dev/test/prod)"
  type        = bool
}

variable "private_subnet_names" {
  description = "List of strings corresponding to private subnet names"
  type        = list(string)
  default     = []
}

variable "public_subnet_names" {
  description = "List of strings corresponding to public subnet names"
  type        = list(string)
  default     = []
}

variable "environment" {
  description = "Environment to be deployed to"
  type        = string
}
