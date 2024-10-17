module "loop_api_gateway_v2" {
  source         = "./modules/api-gateway"
  aws_account_id = data.aws_caller_identity.current.id
  aws_region     = data.aws_region.current.name
  deployment_triggers = [
    jsonencode(module.cloud_file_download_loop_v2),
    jsonencode(module.cloud_file_upload_loop_v2),
    jsonencode(module.cloud_pre_signed_upload_url_loop_v2),
  ]
  domain_name         = module.route53_configuration_loop_api.domain_name_id
  environment         = var.environment
  gateway_name        = "loop_v2"
  acm_certificate_arn = module.route53_configuration_loop_api.acm_certificate_cert_arn
  api_mapping_key     = "v2"
}

module "loop_authorizer_v2_external_http" {
  source                    = "./modules/http-authorizer"
  api_gateway_execution_arn = module.loop_api_gateway_v2.api_gateway_execution_arn
  apigateway_id             = module.loop_api_gateway_v2.api_gateway_id
  authorizer_name           = "sonar_global_external_http_authorizer_v2"
  credentials_role_arn      = module.global.external_authorizer_credentials_role_arn
  identity_sources          = ["$request.header.Authorization"]
  lambda_function_name      = module.global.external_authorizer_lambda_function_name
  lambda_invoke_arn         = module.global.external_authorizer_lambda_invoke_arn
  statement_id              = "AllowExecutionFromAPIGatewayForExternalHTTPAuthorizerV2"
}

module "loop_authorizer_v2_querystring" {
  source                    = "./modules/http-authorizer"
  api_gateway_execution_arn = module.loop_api_gateway_v2.api_gateway_execution_arn
  apigateway_id             = module.loop_api_gateway_v2.api_gateway_id
  authorizer_name           = "sonar_global_external_querystring_authorizer_v2"
  credentials_role_arn      = module.global.external_authorizer_credentials_role_arn
  identity_sources          = ["$request.querystring.authorization"]
  lambda_function_name      = module.global.external_authorizer_lambda_function_name
  lambda_invoke_arn         = module.global.external_authorizer_lambda_invoke_arn
  statement_id              = "AllowExecutionFromAPIGatewayForExternalQuerystringAuthorizerV2"
}