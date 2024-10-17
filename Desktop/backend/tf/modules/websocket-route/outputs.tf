# These are used in CI/CD to calculate SHAs
output "route" {
  value = aws_apigatewayv2_route.route
}

output "integration" {
  value = aws_apigatewayv2_integration.integration
}
