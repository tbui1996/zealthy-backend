variable "environment" {
  type        = string
  description = "Current Environment"
}

variable "db_password" {
  description = "RDS root user password"
  type        = string
  sensitive   = true
}

variable "image_version" {
  type        = string
  description = "Commit SHA of the repository"
}

variable "last_stable_production_commit" {
  type        = string
  description = "last stable production commit for versioned lambdas"
  default     = ""
}

variable "last_stable_develop_commit" {
  type        = string
  description = "last stable production commit for versioned lambdas"
  default     = ""
}

variable "last_stable_test_commit" {
  type        = string
  description = "last stable production commit for versioned lambdas"
  default     = ""
}

variable "last_local_commit" {
  type        = string
  description = "last local commit. DO NOT USE."
  default     = ""
}

variable "hosted_zone_id" {
  type        = string
  description = "Hosted zone ID for developer account"
}

variable "domain_name" {
  type        = string
  description = "Environments Domain name"
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

# vpc vars start
variable "private_subnet_names" {
  description = "List of strings corresponding to private subnet names if VPC already exists"
  type        = list(string)
  default     = []
}

variable "public_subnet_names" {
  description = "List of strings corresponding to public subnet names if VPC already exists"
  type        = list(string)
  default     = []
}
# vpc vars end

variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
}

# doppler db variables
variable "doppler_host" {
  description = "Host string to Doppler DB"
  type        = string
  default     = ""
}

variable "doppler_port" {
  description = "Port for Doppler DB"
  type        = number
  default     = 5432
}

variable "doppler_user" {
  description = "User for logging in to Doppler DB"
  type        = string
  default     = ""
}

variable "doppler_pw" {
  description = "PW for logging in to Doppler DB"
  type        = string
  default     = ""
}

variable "doppler_dbname" {
  description = "Name of Doppler DB"
  type        = string
  default     = "external"
}
