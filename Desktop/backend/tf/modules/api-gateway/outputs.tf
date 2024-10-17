output "api_gateway_id" {
  value = aws_apigatewayv2_api.api_gateway_template.id
}

output "api_gateway_execution_arn" {
  value = aws_apigatewayv2_api.api_gateway_template.execution_arn
}