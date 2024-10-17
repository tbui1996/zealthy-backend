module "router_connect" {
  source    = "./modules/websocket-route"
  route_key = "$connect"

  lambda_function = module.router.lambda_connect

  websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn

  requires_auth = true
  authorizer_id = module.global.external_websocket_authorizer_id
}

module "router_disconnect" {
  source    = "./modules/websocket-route"
  route_key = "$disconnect"

  lambda_function = module.router.lambda_disconnect

  websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn
}

module "router_message" {
  source    = "./modules/websocket-route"
  route_key = "message"

  lambda_function = module.router.lambda_message

  websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn
}

module "router_connector" {
  source = "./modules/router-websocket-connector"

  // environment variables
  domain_name   = var.domain_name
  environment   = var.environment
  image_version = var.image_version

  // resource variables
  name             = "router"
  lambda_image_uri = module.ecr_router_lambda.repository_url

  loop_websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id
  loop_websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn

  router_dynamo_arns = module.router.dynamo_arns
  lambda_receive     = module.router.lambda_receive
}

module "router_unconfirmed_route_connect" {
  source    = "./modules/websocket-route"
  route_key = "$connect"

  lambda_function = module.router.lambda_unconfirmed_connect

  websocket_api_id            = aws_apigatewayv2_api.unconfirmed_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.unconfirmed_websocket_api.execution_arn

  requires_auth = true
  authorizer_id = module.global.external_oh_websocket_authorizer_id
}

module "router_unconfirmed_route_disconnect" {
  source    = "./modules/websocket-route"
  route_key = "$disconnect"

  lambda_function = module.router.lambda_unconfirmed_disconnect

  websocket_api_id            = aws_apigatewayv2_api.unconfirmed_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.unconfirmed_websocket_api.execution_arn
}

module "router_broadcast" {
  source        = "./modules/gateway-route"
  route_key     = local.router_internal_http_routes.broadcast
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.router.lambda_broadcast
}

module "router_user_list" {
  source        = "./modules/gateway-route"
  route_key     = local.deprecated_routes.router_user_list
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.router.lambda_user_list
}
