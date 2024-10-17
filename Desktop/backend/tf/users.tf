module "users_connector" {
  source = "./modules/router-websocket-connector"

  // environment variables
  domain_name   = var.domain_name
  environment   = var.environment
  image_version = var.image_version

  // resource variables
  name             = "users"
  lambda_image_uri = module.ecr_router_lambda.repository_url

  loop_websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id
  loop_websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn

  router_dynamo_arns = module.router.dynamo_arns
  lambda_receive     = module.users.lambda_receive
}

module "users_external_sign_in" {
  source        = "./modules/gateway-route"
  route_key     = local.users_external_http_routes.sign_in
  requires_auth = true

  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_oh_http_authorizer_id

  lambda_function = module.users.lambda_external_sign_in
}

module "users_external_sign_up" {
  source        = "./modules/gateway-route"
  route_key     = local.users_external_http_routes.sign_up
  requires_auth = true

  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_oh_http_authorizer_id

  lambda_function = module.users.lambda_external_sign_up
}

module "users_external_refresh" {
  source        = "./modules/gateway-route"
  route_key     = local.users_external_http_routes.refresh
  requires_auth = true

  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_oh_http_authorizer_id

  lambda_function = module.users.lambda_external_refresh
}

module "users_user_list" {
  source        = "./modules/gateway-route"
  route_key     = local.users_internal_http_routes.user_list
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.users.lambda_user_list
}

module "users_revoke_access" {
  source        = "./modules/gateway-route"
  route_key     = local.users_internal_http_routes.revoke_access
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.users.lambda_revoke_access
}

module "users_get_organizations" {
  source        = "./modules/gateway-route"
  route_key     = local.users_internal_http_routes.get_organizations
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.users.lambda_get_organizations
}

module "users_update_user" {
  source        = "./modules/gateway-route"
  route_key     = local.users_internal_http_routes.update_user
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.users.lambda_update_user
}


module "users_create_organizations" {
  source        = "./modules/gateway-route"
  route_key     = local.users_internal_http_routes.create_organizations
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.users.lambda_create_organizations
}