resource "aws_apigatewayv2_authorizer" "apigatewayv2_authorizer" {
  name                              = var.authorizer_name
  api_id                            = var.apigateway_id
  authorizer_type                   = "REQUEST"
  authorizer_credentials_arn        = var.credentials_role_arn
  authorizer_uri                    = var.lambda_invoke_arn
  identity_sources                  = var.identity_sources
  authorizer_payload_format_version = "1.0"
}

resource "aws_lambda_permission" "permission_apigatewayv2_authorizer" {
  statement_id  = var.statement_id
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.api_gateway_execution_arn}/*/*"
}


