# These are used in CI/CD to calculate SHAs
output "route" {
  value = aws_apigatewayv2_route.http_lambda
}

output "integration" {
  value = aws_apigatewayv2_integration.http_lambda_integration
}
