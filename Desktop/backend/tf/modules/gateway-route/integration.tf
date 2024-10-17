resource "aws_lambda_permission" "apigw_http_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = var.source_arn
}

resource "aws_apigatewayv2_integration" "http_lambda_integration" {
  api_id                 = var.api_id
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  integration_method     = "POST"
  integration_uri        = var.lambda_function.arn
  credentials_arn        = var.credentials_arn != "" ? var.credentials_arn : aws_iam_role.api_gateway_invocation_role.arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "http_lambda" {
  api_id    = var.api_id
  route_key = var.route_key
  target    = "integrations/${aws_apigatewayv2_integration.http_lambda_integration.id}"

  authorization_type = var.requires_auth == true ? "CUSTOM" : "NONE"
  authorizer_id      = var.requires_auth == true ? var.authorizer_id : null
}
