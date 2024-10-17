resource "aws_kms_key" "sonar_group_policy_dynamodb_kms_key" {
  description             = "KMS key for SonarGroupPolicy DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "sonar_group_policy" {
  name         = "SonarGroupPolicy"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "group"

  attribute {
    name = "group"
    type = "S"
  }

  tags = {
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.sonar_group_policy_dynamodb_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}
