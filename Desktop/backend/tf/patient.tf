module "patient_list_appointments" {
  source        = "./modules/gateway-route"
  route_key     = local.appointments_internal_http_routes.list_appointments
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_list_appointments
}

module "patient_create_appointments" {
  source        = "./modules/gateway-route"
  route_key     = local.appointments_internal_http_routes.create_appointments
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_create_appointments
}

module "patient_edit_appointments" {
  source        = "./modules/gateway-route"
  route_key     = local.appointments_internal_http_routes.edit_appointments
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_edit_appointments
}

module "patient_delete_appointments" {
  source        = "./modules/gateway-route"
  route_key     = local.appointments_internal_http_routes.delete_appointments
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_delete_appointments
}

module "patient_list_patients" {
  source        = "./modules/gateway-route"
  route_key     = local.patients_internal_http_routes.list_patients
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_list_patients
}

module "patient_create_patients" {
  source        = "./modules/gateway-route"
  route_key     = local.patients_internal_http_routes.create_patients
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_create_patients
}

module "patient_patch_patients" {
  source        = "./modules/gateway-route"
  route_key     = local.patients_internal_http_routes.patch_patients
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_patch_patients
}

module "patient_list_agency_providers" {
  source        = "./modules/gateway-route"
  route_key     = local.agency_providers_internal_http_routes.list_agency_providers
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_list_agency_providers
}

module "patient_create_agency_providers" {
  source        = "./modules/gateway-route"
  route_key     = local.agency_providers_internal_http_routes.create_agency_providers
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_create_agency_providers
}

module "patient_edit_agency_providers" {
  source        = "./modules/gateway-route"
  route_key     = local.agency_providers_internal_http_routes.edit_agency_providers
  requires_auth = true

  api_id     = aws_apigatewayv2_api.gateway.id
  source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/*/*"

  authorizer_id = module.global.internal_http_authorizer_id

  lambda_function = module.patient.lambda_edit_agency_providers
}

