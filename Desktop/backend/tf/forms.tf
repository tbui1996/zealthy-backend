/*
 * Connects an HTTP endpoint to the corresponding Lambda
 */
module "forms_count" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.count
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_count
}

module "forms_create" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.create
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_create
}

module "forms_get" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.get
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_get
}

module "forms_list" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.list
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_list
}

module "forms_connector" {
  source = "./modules/router-websocket-connector"

  // environment variables
  domain_name   = var.domain_name
  environment   = var.environment
  image_version = var.image_version

  // resource variables
  name             = "forms"
  lambda_image_uri = module.ecr_router_lambda.repository_url

  loop_websocket_api_execution_arn = aws_apigatewayv2_api.loop_websocket_api.execution_arn
  loop_websocket_api_id            = aws_apigatewayv2_api.loop_websocket_api.id

  router_dynamo_arns = module.router.dynamo_arns
  lambda_receive     = module.forms.lambda_receive
}

module "forms_response" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.response
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_response
}

module "forms_send" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.send
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_send
}

module "forms_delete" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.delete
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_delete
}

module "forms_edit" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.edit
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_edit
}

module "forms_close" {
  source        = "./modules/gateway-route"
  route_key     = local.forms_internal_http_routes.close
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.forms.lambda_close
}
