# DynamoDB
resource "aws_kms_key" "websocket_connections_kms_key" {
  description             = "KSM key for SonarWebsocketConnections DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "websocket_connections" {
  name         = "SonarWebsocketConnections"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserID"
  range_key    = "ConnectionId"

  # Enable streaming for Sonar Gateway
  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "UserID"
    type = "S"
  }

  attribute {
    name = "ConnectionId"
    type = "S"
  }

  tags = {
    Name        = "sonar_websocket_connections_db"
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.websocket_connections_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_kms_key" "pending_messages_kms_key" {
  description             = "KSM key for SonarPendingMessages DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "pending_messages" {
  name         = "SonarPendingMessages"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserID"
  range_key    = "CreatedTimestamp"

  attribute {
    name = "UserID"
    type = "S"
  }
  attribute {
    name = "CreatedTimestamp"
    type = "N"
  }

  ttl {
    attribute_name = "DeleteTimestamp"
    enabled        = true
  }

  tags = {
    Name        = "sonar_websocket_connections_db"
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.pending_messages_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_kms_key" "websocket_connections_internal_kms_key" {
  description             = "KSM key for SonarInternalWebsocketConnections DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "websocket_connections_internal" {
  name         = "SonarInternalWebsocketConnections"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserID"
  range_key    = "ConnectionId"

  # Enable streaming for Sonar Gateway
  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "UserID"
    type = "S"
  }

  attribute {
    name = "ConnectionId"
    type = "S"
  }

  attribute {
    name = "CognitoGroup"
    type = "S"
  }

  tags = {
    Name        = "sonar_websocket_connections_db"
    Environment = var.environment
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.websocket_connections_internal_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }

  global_secondary_index {
    name               = "UserGroupIndex"
    hash_key           = "CognitoGroup"
    range_key          = "UserID"
    projection_type    = "INCLUDE"
    non_key_attributes = ["ConnectionId"]
  }
}

resource "aws_kms_key" "unconfirmed_websocket_connections_kms_key" {
  description             = "KMS key for SonarUnconfirmedWebsocketConnections DynamoDB table"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_dynamodb_table" "unconfirmed_websocket_connections" {
  name         = "SonarUnconfirmedWebsocketConnections"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Email"
  range_key    = "ConnectionId"

  # Enable streaming for Sonar Gateway
  stream_enabled   = true
  stream_view_type = "KEYS_ONLY"

  attribute {
    name = "Email"
    type = "S"
  }

  attribute {
    name = "ConnectionId"
    type = "S"
  }

  tags = {
    Name        = "sonar_unconfirmed_websocket_connections_db"
    Environment = var.environment
  }
  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.unconfirmed_websocket_connections_kms_key.arn
  }

  point_in_time_recovery {
    enabled = true
  }
}
