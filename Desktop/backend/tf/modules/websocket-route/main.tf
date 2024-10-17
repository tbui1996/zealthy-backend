resource "aws_apigatewayv2_route" "route" {
  api_id             = var.websocket_api_id
  route_key          = var.route_key
  target             = "integrations/${aws_apigatewayv2_integration.integration.id}"
  authorization_type = var.requires_auth == true ? "CUSTOM" : "NONE"
  authorizer_id      = var.requires_auth == true ? var.authorizer_id : null
}

resource "aws_apigatewayv2_integration" "integration" {
  api_id                    = var.websocket_api_id
  integration_type          = "AWS_PROXY"
  integration_method        = "POST"
  integration_uri           = var.lambda_function.invoke_arn
  content_handling_strategy = "CONVERT_TO_TEXT"
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.websocket_api_execution_arn}/*/*"
}