locals {
  http_gateway_logs_name               = "sonar_${var.gateway_name}_gateway_http_api_logs"
  http_gateway_logs_kms_key_alias_name = "alias/http_${var.gateway_name}_gateway_logs_kms_key"
}