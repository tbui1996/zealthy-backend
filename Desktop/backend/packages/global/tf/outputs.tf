output "internal_websocket_authorizer_id" {
  value = aws_apigatewayv2_authorizer.internal_websocket_authorizer.id
}

output "internal_http_authorizer_id" {
  value = aws_apigatewayv2_authorizer.internal_http_authorizer.id
}

output "external_websocket_authorizer_id" {
  value = aws_apigatewayv2_authorizer.external_websocket_authorizer.id
}

output "external_http_authorizer_id" {
  value = aws_apigatewayv2_authorizer.external_http_authorizer.id
}

output "external_oh_websocket_authorizer_id" {
  value = aws_apigatewayv2_authorizer.external_oh_websocket_authorizer.id
}

output "external_oh_http_authorizer_id" {
  value = aws_apigatewayv2_authorizer.external_oh_http_authorizer.id
}

output "lambda_connect" {
  value = aws_lambda_function.internal_websocket_connect
}

output "lambda_disconnect" {
  value = aws_lambda_function.internal_websocket_disconnect
}

output "task_is_internal_user_online" {
  value = aws_lambda_function.global_task_is_internal_user_online
}

output "external_querystring_authorizer_id" {
  value = aws_apigatewayv2_authorizer.external_querystring_authorizer.id
}

output "internal_querystring_authorizer_id" {
  value = aws_apigatewayv2_authorizer.internal_querystring_authorizer.id
}

output "external_authorizer_credentials_role_arn" {
  value = aws_iam_role.external_authorizer_credentials_role.arn
}

output "external_authorizer_lambda_invoke_arn" {
  value = aws_lambda_function.external_authorizer_lambda.invoke_arn
}

output "external_authorizer_lambda_function_name" {
  value = aws_lambda_function.external_authorizer_lambda.function_name
}