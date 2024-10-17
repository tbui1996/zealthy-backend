module "global_connect_websocket_route" {
  source    = "./modules/websocket-route"
  route_key = "$connect"

  lambda_function = module.global.lambda_connect

  websocket_api_id            = aws_apigatewayv2_api.webapp_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.webapp_websocket_api.execution_arn

  requires_auth = true
  authorizer_id = module.global.internal_websocket_authorizer_id
}

module "global_disconnect_websocket_route" {
  source    = "./modules/websocket-route"
  route_key = "$disconnect"

  lambda_function = module.global.lambda_disconnect

  websocket_api_id            = aws_apigatewayv2_api.webapp_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.webapp_websocket_api.execution_arn
}
