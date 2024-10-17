resource "aws_kms_key" "messages_kms_key" {
  description             = "KMS key for Messages DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "messages" {
  name         = "SonarMessages"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "SessionID"
  range_key    = "CreatedTimestamp"

  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "SessionID"
    type = "S"
  }

  attribute {
    name = "CreatedTimestamp"
    type = "N"
  }

  tags = {
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.messages_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_kms_key" "feedback_kms_key" {
  description             = "KMS key for Feedback DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "feedback" {
  name         = "SonarFeedback"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "email"
  range_key    = "createdTimestamp"

  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "email"
    type = "S"
  }

  attribute {
    name = "createdTimestamp"
    type = "N"
  }

  tags = {
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.feedback_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_kms_key" "offline_messages_kms_key" {
  description             = "KMS key for Offline Messages DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "offline_messages" {
  name         = "OfflineMessageNotifications"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserID"

  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "UserID"
    type = "S"
  }

  tags = {
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.offline_messages_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}