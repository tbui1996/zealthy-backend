resource "aws_kms_key" "service_send_kms_key" {
  description             = "KMS key for sonar-service-${var.name}-send SQS queue"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_sqs_queue" "service_send" {
  name                      = "sonar-service-${var.name}-send"
  message_retention_seconds = 60

  kms_master_key_id                 = aws_kms_key.service_send_kms_key.key_id
  kms_data_key_reuse_period_seconds = 300

  tags = {
    Environment = var.environment
  }
}

resource "aws_kms_key" "service_receive_kms_key" {
  description             = "KMS key for sonar-service-${var.name}-receive SQS queue"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_sqs_queue" "service_receive" {
  name                      = "sonar-service-${var.name}-receive"
  message_retention_seconds = 60

  kms_master_key_id                 = aws_kms_key.service_receive_kms_key.key_id
  kms_data_key_reuse_period_seconds = 300

  tags = {
    Environment = var.environment
  }
}

resource "aws_lambda_event_source_mapping" "dispatch_send" {
  # TODO: Address in SONAR-298
  batch_size       = 1
  event_source_arn = aws_sqs_queue.service_send.arn
  function_name    = aws_lambda_function.route_service_forward.function_name
}

resource "aws_lambda_event_source_mapping" "dispatch_receive" {
  # TODO: Address in SONAR-298
  batch_size       = 1
  event_source_arn = aws_sqs_queue.service_receive.arn
  function_name    = var.lambda_receive.function_name
}
