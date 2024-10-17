
module "feature_flags_create_flag" {
  source        = "./modules/gateway-route"
  route_key     = local.feature_flags_internal_http_routes.create_flag
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.feature_flags.lambda_create_flag
}

module "feature_flags_patch_flag" {
  source        = "./modules/gateway-route"
  route_key     = local.feature_flags_internal_http_routes.patch_flag
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.feature_flags.lambda_patch_flag
}


module "feature_flags_list_flags" {
  source        = "./modules/gateway-route"
  route_key     = local.feature_flags_internal_http_routes.list_flags
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.feature_flags.lambda_list_flags
}

module "feature_flags_evaluate" {
  source        = "./modules/gateway-route"
  route_key     = local.feature_flags_internal_http_routes.evaluate
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.feature_flags.lambda_evaluate
}

module "feature_flags_loop_evaluate" {
  source        = "./modules/gateway-route"
  route_key     = local.feature_flags_external_http_routes.evaluate
  requires_auth = true

  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.feature_flags.lambda_loop_evaluate
}

module "feature_flags_delete_flag" {
  source        = "./modules/gateway-route"
  route_key     = local.feature_flags_internal_http_routes.delete_flag
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.feature_flags.lambda_delete_flag
}
