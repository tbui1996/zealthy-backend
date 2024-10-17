module "cloud_file_upload_loop" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_external_http_routes.file_upload
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_http_authorizer_id

  lambda_function = module.cloud.lambda_file_upload_loop
}

module "cloud_file_upload_loop_v2" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_external_v2_http_routes.file_upload
  requires_auth = true

  // module router variables
  api_id     = module.loop_api_gateway_v2.api_gateway_id
  source_arn = "${module.loop_api_gateway_v2.api_gateway_execution_arn}/*/*"

  // module global variables
  authorizer_id = module.loop_authorizer_v2_external_http.external_apigatewayv2_authorizer_id

  lambda_function = module.cloud.lambda_file_upload_loop_v2
}

module "cloud_file_download_loop" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_external_http_routes.file_download
  requires_auth = true

  // module router variables
  api_id     = aws_apigatewayv2_api.loop_gateway.id
  source_arn = "${aws_apigatewayv2_api.loop_gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.external_querystring_authorizer_id

  lambda_function = module.cloud.lambda_file_download_loop
}

module "cloud_file_download_loop_v2" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_external_v2_http_routes.file_download
  requires_auth = true

  // module router variables
  api_id     = module.loop_api_gateway_v2.api_gateway_id
  source_arn = "${module.loop_api_gateway_v2.api_gateway_execution_arn}/*/*"

  // module global variables
  authorizer_id = module.loop_authorizer_v2_querystring.external_apigatewayv2_authorizer_id

  lambda_function = module.cloud.lambda_file_download_loop_v2
}

// Internal/web app routes
module "cloud_file_download_web" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_internal_http_routes.file_download
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_querystring_authorizer_id

  lambda_function = module.cloud.lambda_file_download_web
}

module "cloud_file_upload_web" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_internal_http_routes.file_upload
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.cloud.lambda_file_upload_web
}

module "cloud_get_file" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_internal_http_routes.get_file
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.cloud.lambda_get_file
}

module "cloud_associate_file" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_internal_http_routes.associate_file
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.cloud.lambda_associate_file
}

module "cloud_delete_file" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_internal_http_routes.delete_file
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.cloud.lambda_delete_file
}

module "cloud_pre_signed_upload_url_web" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_internal_http_routes.pre_signed_upload_url
  requires_auth = true

  // resource variables
  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  // module global variables
  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.cloud.lambda_pre_signed_upload_url_web
}

module "cloud_pre_signed_upload_url_loop_v2" {
  source        = "./modules/gateway-route"
  route_key     = local.cloud_external_v2_http_routes.pre_signed_upload_url
  requires_auth = true

  // module router variables
  api_id     = module.loop_api_gateway_v2.api_gateway_id
  source_arn = "${module.loop_api_gateway_v2.api_gateway_execution_arn}/*/*"

  // module global variables
  authorizer_id = module.loop_authorizer_v2_external_http.external_apigatewayv2_authorizer_id

  lambda_function = module.cloud.lambda_pre_signed_upload_url_loop
}
