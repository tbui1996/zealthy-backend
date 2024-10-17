output "dynamo_arns" {
  description = "Dynamo arns for tables and keys"

  value = {
    websocket_connections                     = aws_dynamodb_table.websocket_connections.arn
    websocket_connections_kms_key             = aws_kms_key.websocket_connections_kms_key.arn
    pending_messages                          = aws_dynamodb_table.pending_messages.arn
    pending_messages_kms_key                  = aws_kms_key.pending_messages_kms_key.arn
    websocket_connections_internal            = aws_dynamodb_table.websocket_connections_internal.arn
    websocket_connections_internal_kms_key    = aws_kms_key.websocket_connections_internal_kms_key.arn
    unconfirmed_websocket_connections         = aws_dynamodb_table.unconfirmed_websocket_connections.arn
    unconfirmed_websocket_connections_kms_key = aws_kms_key.unconfirmed_websocket_connections_kms_key.arn
  }
}

output "lambda_broadcast" {
  value = aws_lambda_function.api_broadcast_lambda
}

output "lambda_user_list" {
  value = aws_lambda_function.api_users_list
}

output "lambda_connect" {
  value = aws_lambda_function.sonar_service_router_connect
}

output "lambda_disconnect" {
  value = aws_lambda_function.sonar_service_router_disconnect
}

output "lambda_message" {
  value = aws_lambda_function.sonar_service_router_message
}

output "lambda_unconfirmed_connect" {
  value = aws_lambda_function.sonar_service_router_unconfirmed_connect
}

output "lambda_unconfirmed_disconnect" {
  value = aws_lambda_function.sonar_service_unconfirmed_disconnect
}

output "lambda_receive" {
  value = aws_lambda_function.sonar_service_router_receive
}

output "task_is_external_user_online" {
  value = aws_lambda_function.sonar_service_router_task_is_external_user_online
}