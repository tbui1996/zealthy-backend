module "router" {
  source = "../packages/router/tf"

  domain_name    = var.domain_name
  environment    = var.environment
  hosted_zone_id = var.hosted_zone_id
  image_version  = var.image_version

  lambda_image_uri = module.ecr_router_lambda.repository_url

  name = "router"

  db_name     = aws_rds_cluster.sonar_service_rds_cluster.database_name
  db_host     = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_port     = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_password = var.db_password

  # RDS Configuration
  vpc_id             = module.vpc.vpc_id
  private_subnets    = module.vpc.private_subnets
  rds_security_group = aws_security_group.rds_security_group.id

  receive_queue_arn     = module.router_connector.receive_queue_arn
  receive_queue_kms_arn = module.router_connector.receive_queue_kms_arn
  send_queue_arn        = module.router_connector.send_queue_arn
  send_queue_kms_arn    = module.router_connector.send_queue_kms_arn

  loop_websocket_api_id = aws_apigatewayv2_api.loop_websocket_api.id

  event_bus_arn = aws_cloudwatch_event_bus.service_events.arn
}

module "forms" {
  source = "../packages/forms/tf"

  domain_name   = var.domain_name
  image_version = var.image_version
  environment   = var.environment

  lambda_image_uri = module.ecr_forms_lambda.repository_url

  # Allow access to DB
  db_host      = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_password  = var.db_password
  db_port      = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username  = aws_rds_cluster.sonar_service_rds_cluster.master_username
  form_db_name = aws_rds_cluster.sonar_service_rds_cluster.database_name

  # RDS Configuration
  vpc_id             = module.vpc.vpc_id
  private_subnets    = module.vpc.private_subnets
  rds_security_group = aws_security_group.rds_security_group.id

  external_security_group = aws_security_group.external_service_security_group.id

  receive_queue_arn     = module.forms_connector.receive_queue_arn
  receive_queue_kms_arn = module.forms_connector.receive_queue_kms_arn
  send_queue_arn        = module.forms_connector.send_queue_arn
  send_queue_kms_arn    = module.forms_connector.send_queue_kms_arn
}

module "support" {
  source = "../packages/support/tf"

  environment   = var.environment
  image_version = var.image_version
  domain_name   = var.domain_name

  lambda_image_uri = module.ecr_support_lambda.repository_url

  apigateway_id         = aws_apigatewayv2_api.gateway.id
  apigateway_stage_name = aws_apigatewayv2_stage.v1.name

  webapp_websocket_api_id = aws_apigatewayv2_api.webapp_websocket_api.id

  aws_apigatewayv2_loop_gateway_id         = aws_apigatewayv2_api.loop_gateway.id
  aws_apigatewayv2_loop_gateway_stage_name = aws_apigatewayv2_stage.loop_v1.name

  loop_websocket_api_id = aws_apigatewayv2_api.loop_websocket_api.id

  router_dynamo_arns = module.router.dynamo_arns

  receive_queue_arn     = module.support_connector.receive_queue_arn
  receive_queue_kms_arn = module.support_connector.receive_queue_kms_arn
  send_queue_arn        = module.support_connector.send_queue_arn
  send_queue_kms_arn    = module.support_connector.send_queue_kms_arn

  # RDS Configuration
  vpc_id             = module.vpc.vpc_id
  private_subnets    = module.vpc.private_subnets
  rds_security_group = aws_security_group.rds_security_group.id

  db_host     = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_password = var.db_password
  db_port     = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_name     = aws_rds_cluster.sonar_service_rds_cluster.database_name

  external_security_group = aws_security_group.external_service_security_group.id

  # Email Configuration
  email_template_feedback         = aws_ses_template.feedback.name
  email_template_new_message_open = aws_ses_template.new_message_open.name
  configuration_set_name          = aws_ses_configuration_set.default_event_publishing.name
  email_identity                  = aws_ses_domain_identity.default.arn

  event_bus_arn = aws_cloudwatch_event_bus.service_events.arn
}

module "global" {
  source           = "../packages/global/tf"
  name             = "global"
  lambda_image_uri = module.ecr_global_lambda.repository_url
  image_version    = var.image_version

  // web api gw
  webapp_websocket_api_id            = aws_apigatewayv2_api.webapp_websocket_api.id
  webapp_websocket_api_execution_arn = aws_apigatewayv2_api.webapp_websocket_api.execution_arn
  webapp_http_api_id                 = aws_apigatewayv2_api.gateway.id
  webapp_http_api_execution_arn      = aws_apigatewayv2_api.gateway.execution_arn

  // router api gw
  loop_websocket_api_id                        = aws_apigatewayv2_api.loop_websocket_api.id
  loop_websocket_api_execution_arn             = aws_apigatewayv2_api.loop_websocket_api.execution_arn
  loop_http_api_id                             = aws_apigatewayv2_api.loop_gateway.id
  loop_http_api_execution_arn                  = aws_apigatewayv2_api.loop_gateway.execution_arn
  loop_unconfirmed_websocket_api_id            = aws_apigatewayv2_api.unconfirmed_websocket_api.id
  loop_unconfirmed_websocket_api_execution_arn = aws_apigatewayv2_api.unconfirmed_websocket_api.execution_arn

  router_dynamo_arns                   = module.router.dynamo_arns
  aws_cognito_user_pool_internals_jwks = module.users.aws_cognito_user_pool_internals_jwks
  aws_cognito_user_pool_externals_jwks = module.users.aws_cognito_user_pool_externals_jwks
  environment                          = var.environment

  group_names  = toset(concat(keys(local.internal_groups), keys(local.external_groups)))
  group_routes = merge(local.internal_groups, local.external_groups)

  // grafana
  vpc_id              = module.vpc.vpc_id
  vpc_public_subnets  = module.vpc.public_subnets
  vpc_private_subnets = module.vpc.private_subnets
  domain_name         = var.domain_name
  db_host             = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_password         = var.db_password
  db_port             = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username         = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_name             = aws_rds_cluster.sonar_service_rds_cluster.database_name
  aws_account_id      = var.aws_account_id
  ecs_cluster_name    = aws_ecs_cluster.sonar_cluster.name
  hosted_zone_id      = var.hosted_zone_id

  event_bus_arn = aws_cloudwatch_event_bus.service_events.arn
}

module "users" {
  source        = "../packages/users/tf"
  domain_name   = var.domain_name
  environment   = var.environment
  image_version = var.image_version

  lambda_image_uri = module.ecr_users_lambda.repository_url

  router_dynamo_arns = module.router.dynamo_arns

  okta_org_name  = var.okta_org_name
  okta_base_url  = var.okta_base_url
  okta_api_token = var.okta_api_token
  callback_urls  = var.callback_urls
  logout_urls    = var.logout_urls

  loop_unconfirmed_websocket_api_id = aws_apigatewayv2_api.unconfirmed_websocket_api.id
  loop_websocket_api_id             = aws_apigatewayv2_api.loop_websocket_api.id

  # normalized group names for cognito
  internal_group_names = toset(keys(local.internal_groups))
  external_group_names = toset(keys(local.external_groups))

  # user_id and username are used for developer accounts only
  okta_user_id  = var.okta_user_id
  okta_username = var.okta_username
  live_env      = local.live_env

  receive_queue_arn     = module.users_connector.receive_queue_arn
  receive_queue_kms_arn = module.users_connector.receive_queue_kms_arn
  send_queue_arn        = module.users_connector.send_queue_arn
  send_queue_kms_arn    = module.users_connector.send_queue_kms_arn

  # RDS Configuration
  vpc_id                  = module.vpc.vpc_id
  private_subnets         = module.vpc.private_subnets
  rds_security_group      = aws_security_group.rds_security_group.id
  external_security_group = aws_security_group.external_service_security_group.id
  db_host                 = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_password             = var.db_password
  db_port                 = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username             = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_name                 = aws_rds_cluster.sonar_service_rds_cluster.database_name
}

module "cloud" {
  source        = "../packages/cloud/tf"
  domain_name   = var.domain_name
  image_version = var.image_version
  environment   = var.environment

  lambda_image_uri = module.ecr_cloud_lambda.repository_url

  # Allow access to RDS main database
  db_password = var.db_password
  db_host     = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_port     = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_name     = aws_rds_cluster.sonar_service_rds_cluster.database_name

  # RDS Configuration
  vpc_id             = module.vpc.vpc_id
  private_subnets    = module.vpc.private_subnets
  rds_security_group = aws_security_group.rds_security_group.id

  # S3 VPC Configuration
  external_security_group = aws_security_group.external_service_security_group.id

  # Local developer should not have to worry about the backwards compatibility. If they do, it will be documented on how to test against multiple different versions
  last_production_commit = var.image_version
  # uncomment the line below (and comment out line above) only when steps to introduce backwards compatibility have been taken
  # last_production_commit = local.live_env ? lookup(local.live_envs_last_stable_commit, var.environment, var.last_local_commit) : var.image_version

  # uncomment the line below to test backwards compatibility locally
  # last_production_commit = lookup(local.live_envs_last_stable_commit, var.environment, var.last_local_commit)

  # Doppler DB Config
  doppler_host   = var.doppler_host
  doppler_port   = var.doppler_port
  doppler_user   = var.doppler_user
  doppler_pw     = var.doppler_pw
  doppler_dbname = var.doppler_dbname
  live_env       = local.live_env
}

module "feature_flags" {
  source = "../packages/feature_flags/tf"

  image_version = var.image_version
  environment   = var.environment
  domain_name   = var.domain_name

  lambda_image_uri = module.ecr_feature_flags_lambda.repository_url

  apigateway_id         = aws_apigatewayv2_api.gateway.id
  apigateway_stage_name = aws_apigatewayv2_stage.v1.name

  aws_apigatewayv2_loop_gateway_id         = aws_apigatewayv2_api.loop_gateway.id
  aws_apigatewayv2_loop_gateway_stage_name = aws_apigatewayv2_stage.loop_v1.name

  # RDS Configuration
  vpc_id                  = module.vpc.vpc_id
  private_subnets         = module.vpc.private_subnets
  rds_security_group      = aws_security_group.rds_security_group.id
  external_security_group = aws_security_group.external_service_security_group.id
  db_host                 = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_read_only_host       = aws_rds_cluster.sonar_service_rds_cluster.reader_endpoint
  db_password             = var.db_password
  db_port                 = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username             = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_name                 = aws_rds_cluster.sonar_service_rds_cluster.database_name
}

module "patient" {
  source        = "../packages/patient/tf"
  domain_name   = var.domain_name
  image_version = var.image_version
  environment   = var.environment

  lambda_image_uri = module.ecr_patient_lambda.repository_url

  apigateway_id         = aws_apigatewayv2_api.gateway.id
  apigateway_stage_name = aws_apigatewayv2_stage.v1.name

  # RDS Configuration
  vpc_id             = module.vpc.vpc_id
  private_subnets    = module.vpc.private_subnets
  rds_security_group = aws_security_group.rds_security_group.id

  # Allow access to RDS main database
  db_password = var.db_password
  db_host     = aws_rds_cluster.sonar_service_rds_cluster.endpoint
  db_port     = aws_rds_cluster.sonar_service_rds_cluster.port
  db_username = aws_rds_cluster.sonar_service_rds_cluster.master_username
  db_name     = aws_rds_cluster.sonar_service_rds_cluster.database_name

  # Doppler DB Config
  doppler_host   = var.doppler_host
  doppler_port   = var.doppler_port
  doppler_user   = var.doppler_user
  doppler_pw     = var.doppler_pw
  doppler_dbname = var.doppler_dbname
  live_env       = local.live_env
}
