module "support_connector" {
  source = "./modules/router-websocket-connector"

  // environment variables
  domain_name   = var.domain_name
  environment   = var.environment
  image_version = var.image_version

  // resource variables
  name             = "support"
  lambda_image_uri = module.ecr_router_lambda.repository_url

  loop_websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id
  loop_websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn

  router_dynamo_arns = module.router.dynamo_arns
  lambda_receive     = module.support.lambda_receive
}

module "support_chat_session" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.chat_session
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_chat_session
}

module "support_loop_pending_chat_session" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.chat_session
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_pending_chat_session_create
}

module "support_pending_chat_sessions_get" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.pending_chat_sessions
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_pending_chat_sessions_get
}

module "support_assign_pending_chat_session" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.assign_pending_chat_session
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_assign_pending_chat_session
}

module "support_chat_messages_get" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.get_session_messages
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_chat_messages_get
}

module "support_chat_sessions_get" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.get_sessions
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_chat_sessions_get
}

module "support_chat_sessions_update_open" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.update_chat_session
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_chat_session_update_open
}

module "support_websocket_route" {
  source    = "./modules/websocket-route"
  route_key = "support"

  lambda_function = module.support.lambda_send

  websocket_api_id            = aws_apigatewayv2_api.webapp_websocket_api.id
  websocket_api_execution_arn = aws_apigatewayv2_api.webapp_websocket_api.execution_arn
}


module "support_loop_chat_messages_get" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.get_session_messages
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_loop_chat_messages_get
}

module "support_loop_chat_sessions_get" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.get_user_sessions
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_loop_chat_sessions_get
}

module "support_loop_send_route" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.send_message
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_loop_send
}

module "support_update_chat_notes" {
  source        = "./modules/gateway-route"
  route_key     = local.support_internal_http_routes.update_chat_notes
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.support.lambda_update_chat_notes
}

module "support_submit_feedback" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.submit_feedback
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_submit_feedback
}

module "support_chat_session_star" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.chat_session_star
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_chat_session_star
}

module "support_loop_online_internal_users" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.online_internal_users
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_loop_online_internal_users
}

module "support_loop_patients_get" {
  source        = "./modules/gateway-route"
  route_key     = local.support_external_http_routes.patients_get
  requires_auth = true

  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.support.lambda_loop_patients_get
}