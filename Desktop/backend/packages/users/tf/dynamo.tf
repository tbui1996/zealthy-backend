resource "aws_kms_key" "domain_whitelist_kms_key" {
  description             = "KMS key for SonarEmailDomainWhitelist DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "domain_whitelist" {
  name         = "SonarEmailDomainWhitelist"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "EmailDomain"

  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "EmailDomain"
    type = "S"
  }

  tags = {
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.domain_whitelist_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_dynamodb_table_item" "circulo" {
  table_name = aws_dynamodb_table.domain_whitelist.name
  hash_key   = aws_dynamodb_table.domain_whitelist.hash_key

  item = <<ITEM
{
  "EmailDomain": {"S": "circulohealth.com"}
}
ITEM
}

resource "aws_dynamodb_table_item" "olive" {
  table_name = aws_dynamodb_table.domain_whitelist.name
  hash_key   = aws_dynamodb_table.domain_whitelist.hash_key

  item = <<ITEM
{
  "EmailDomain": {"S": "oliveai.com "}
}
ITEM
}
