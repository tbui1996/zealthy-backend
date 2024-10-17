################################################################################
#############################    WEBSOCKET AUTHORIZER  #########################
resource "aws_apigatewayv2_authorizer" "internal_websocket_authorizer" {
  name                       = "sonar_global_internal_websocket_authorizer"
  api_id                     = var.webapp_websocket_api_id
  authorizer_type            = "REQUEST"
  authorizer_credentials_arn = aws_iam_role.internal_authorizer_credentials_role.arn
  authorizer_uri             = aws_lambda_function.internal_authorizer_lambda.invoke_arn
  identity_sources           = ["route.request.querystring.authorization"]
}

resource "aws_lambda_permission" "internal_websocket_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayInternalWebsocketAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.internal_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.webapp_websocket_api_execution_arn}/*/*"
}

resource "aws_apigatewayv2_authorizer" "external_websocket_authorizer" {
  name                       = "sonar_global_external_websocket_authorizer"
  api_id                     = var.loop_websocket_api_id
  authorizer_type            = "REQUEST"
  authorizer_credentials_arn = aws_iam_role.external_authorizer_credentials_role.arn
  authorizer_uri             = aws_lambda_function.external_authorizer_lambda.invoke_arn
  identity_sources           = ["route.request.header.Authorization"]
}

resource "aws_lambda_permission" "external_websocket_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayExternalWebsocketAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.external_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.loop_websocket_api_execution_arn}/*/*"
}

resource "aws_apigatewayv2_authorizer" "external_oh_websocket_authorizer" {
  name                       = "sonar_global_external_oh_websocket_authorizer"
  api_id                     = var.loop_unconfirmed_websocket_api_id
  authorizer_type            = "REQUEST"
  authorizer_credentials_arn = aws_iam_role.external_oh_authorizer_credentials_role.arn
  authorizer_uri             = aws_lambda_function.external_oh_authorizer_lambda.invoke_arn
  identity_sources           = ["route.request.header.Authorization"]
}

resource "aws_lambda_permission" "external_oh_websocket_authorizer" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.external_oh_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.loop_unconfirmed_websocket_api_execution_arn}/*/*"
}



#############################    WEBSOCKET AUTHORIZER  #########################
################################################################################

################################################################################
#############################    HTTP AUTHORIZER  #########################
resource "aws_apigatewayv2_authorizer" "internal_http_authorizer" {
  name                              = "sonar_global_internal_http_authorizer"
  api_id                            = var.webapp_http_api_id
  identity_sources                  = ["$request.header.Authorization"]
  authorizer_type                   = "REQUEST"
  authorizer_credentials_arn        = aws_iam_role.internal_authorizer_credentials_role.arn
  authorizer_uri                    = aws_lambda_function.internal_authorizer_lambda.invoke_arn
  authorizer_payload_format_version = "1.0"
}

resource "aws_lambda_permission" "internal_http_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayInternalHTTPAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.internal_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.webapp_http_api_execution_arn}/*/*"
}

resource "aws_apigatewayv2_authorizer" "external_http_authorizer" {
  name                              = "sonar_global_external_http_authorizer"
  api_id                            = var.loop_http_api_id
  identity_sources                  = ["$request.header.Authorization"]
  authorizer_type                   = "REQUEST"
  authorizer_credentials_arn        = aws_iam_role.external_authorizer_credentials_role.arn
  authorizer_uri                    = aws_lambda_function.external_authorizer_lambda.invoke_arn
  authorizer_payload_format_version = "1.0"
}

resource "aws_lambda_permission" "external_http_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayExternalHTTPAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.external_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.loop_http_api_execution_arn}/*/*"
}

resource "aws_apigatewayv2_authorizer" "external_oh_http_authorizer" {
  name                              = "sonar_global_external_oh_http_authorizer"
  api_id                            = var.loop_http_api_id
  authorizer_type                   = "REQUEST"
  authorizer_credentials_arn        = aws_iam_role.external_oh_authorizer_credentials_role.arn
  authorizer_uri                    = aws_lambda_function.external_oh_authorizer_lambda.invoke_arn
  identity_sources                  = ["$request.header.Authorization"]
  authorizer_payload_format_version = "1.0"
  authorizer_result_ttl_in_seconds  = 0
}

resource "aws_lambda_permission" "external_oh_http_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayForExternalOHHTTPAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.external_oh_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.loop_http_api_execution_arn}/*/*"
}

#############################    HTTP AUTHORIZER  #########################
################################################################################

################################################################################
############################# QUERYSTRING AUTHORIZER  #########################
resource "aws_apigatewayv2_authorizer" "internal_querystring_authorizer" {
  name                              = "sonar_global_internal_querystring_authorizer"
  api_id                            = var.webapp_http_api_id
  identity_sources                  = ["$request.querystring.authorization"]
  authorizer_type                   = "REQUEST"
  authorizer_credentials_arn        = aws_iam_role.internal_authorizer_credentials_role.arn
  authorizer_uri                    = aws_lambda_function.internal_authorizer_lambda.invoke_arn
  authorizer_payload_format_version = "1.0"
}

resource "aws_lambda_permission" "internal_querystring_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayInternalQuerystringAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.internal_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.webapp_http_api_execution_arn}/*/*"
}

resource "aws_apigatewayv2_authorizer" "external_querystring_authorizer" {
  name                              = "sonar_global_external_querystring_authorizer"
  api_id                            = var.loop_http_api_id
  identity_sources                  = ["$request.querystring.authorization"]
  authorizer_type                   = "REQUEST"
  authorizer_credentials_arn        = aws_iam_role.external_authorizer_credentials_role.arn
  authorizer_uri                    = aws_lambda_function.external_authorizer_lambda.invoke_arn
  authorizer_payload_format_version = "1.0"
}

resource "aws_lambda_permission" "external_querystring_authorizer" {
  statement_id  = "AllowExecutionFromAPIGatewayExternalQuerystringAuthorizer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.external_authorizer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.loop_http_api_execution_arn}/*/*"
}
############################# QUERYSTRING AUTHORIZER  #########################
################################################################################