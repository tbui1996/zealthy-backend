resource "aws_kms_key" "http_gateway_logs_kms_key" {
  description             = "KMS key for sonar-${var.gateway_name}-gateway-http-api-logs CloudWatch log group"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.http_gateway_logs_kms_key_document.json
}

resource "aws_kms_alias" "http_gateway_logs_kms_key" {
  name          = local.http_gateway_logs_kms_key_alias_name
  target_key_id = aws_kms_key.http_gateway_logs_kms_key.id
}