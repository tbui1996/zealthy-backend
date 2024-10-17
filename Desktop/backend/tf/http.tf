# Give CLoudWatch Log Group KMS permission
# https://github.com/hashicorp/terraform-provider-aws/issues/8042
# https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/encrypt-log-data-kms.html
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key#policy
locals {
  http_gateway_logs_name               = "sonar_webapp_gateway_http_api_logs"
  http_gateway_logs_kms_key_alias_name = "alias/gateway_logs_kms_key"
}

resource "aws_kms_key" "gateway_logs_kms_key" {
  description             = "KMS key for sonar_webapp_gateway_http_api_logs CloudWatch log group"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.gateway_logs_kms_key_document.json
}

resource "aws_kms_alias" "gateway_logs_kms_key" {
  name          = local.http_gateway_logs_kms_key_alias_name
  target_key_id = aws_kms_key.gateway_logs_kms_key.id
}

data "aws_iam_policy_document" "gateway_logs_kms_key_document" {
  policy_id = "sonar-webapp-gateway-http-api-logs-kms-key"
  statement {
    sid = "Enable IAM User Permissions"
    principals {
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
      type        = "AWS"
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }

  statement {
    sid = "Enable key to be used only with Cloud Watch Log"
    principals {
      identifiers = ["logs.${data.aws_region.current.name}.amazonaws.com"]
      type        = "Service"
    }
    actions = [
      "kms:Encrypt*",
      "kms:Decrypt*",
      "kms:ReEncryptTo",
      "kms:ReEncryptFrom",
      "kms:GenerateDataKey*",
      "kms:Describe*"
    ]
    resources = ["*"]
    condition {
      test     = "ArnEquals"
      variable = "kms:EncryptionContext:aws:logs:arn"
      values   = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:${local.http_gateway_logs_name}"]
    }
  }
}

resource "aws_cloudwatch_log_group" "gateway_http_logs" {
  name       = local.http_gateway_logs_name
  kms_key_id = aws_kms_key.gateway_logs_kms_key.arn
}

resource "aws_apigatewayv2_api" "gateway" {
  name          = "sonar_webapp_gateway_http_api"
  protocol_type = "HTTP"
  cors_configuration {
    allow_headers = ["Content-Type", "X-Amz-Date", "Authorization", "X-Api-Key", "X-Amz-Security-Token"]
    allow_methods = ["GET", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"]
    allow_origins = ["*"]
  }
}

resource "aws_apigatewayv2_stage" "v1" {
  api_id        = aws_apigatewayv2_api.gateway.id
  name          = var.environment
  deployment_id = aws_apigatewayv2_deployment.http_deploy.id

  default_route_settings {
    logging_level          = "INFO"
    throttling_rate_limit  = 1000
    throttling_burst_limit = 1000
  }

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.gateway_http_logs.arn
    # JSON log format so that inside cloudwatch, the available values can be inferred when querying
    format = jsonencode({
      requestId : "$context.requestId",
      ip : "$context.identity.sourceIp",
      caller : "$context.identity.caller",
      # userID is the userID parsed from the sonar token generated by cognito, and should be included in all api calls
      userID : "$context.authorizer.userID",
      user : "$context.identity.user",
      requestTime : "$context.requestTime",
      httpMethod : "$context.httpMethod",
      resourcePath : "$context.resourcePath",
      path : "$context.path",
      status : "$context.status",
      protocol : "$context.protocol",
      responseLength : "$context.responseLength",
      authorizerError : "$context.authorizer.error",
      errorMessage : "$context.error.message",
      integrationError : "$context.integration.error",
    })
  }
}

resource "aws_apigatewayv2_deployment" "http_deploy" {
  api_id = aws_apigatewayv2_api.gateway.id

  lifecycle {
    create_before_destroy = true
  }

  triggers = {
    redeployment = sha1(join(",", [
      jsonencode(module.cloud_file_download_web),
      jsonencode(module.cloud_file_upload_web),
      jsonencode(module.cloud_get_file),
      jsonencode(module.cloud_associate_file),
      jsonencode(module.cloud_delete_file),
      jsonencode(module.cloud_pre_signed_upload_url_web),

      jsonencode(module.forms_count),
      jsonencode(module.forms_create),
      jsonencode(module.forms_get),
      jsonencode(module.forms_list),
      jsonencode(module.forms_response),
      jsonencode(module.forms_send),
      jsonencode(module.forms_delete),
      jsonencode(module.forms_edit),
      jsonencode(module.forms_close),

      jsonencode(module.router_broadcast),
      jsonencode(module.router_user_list),

      jsonencode(module.support_chat_session),
      jsonencode(module.support_pending_chat_sessions_get),
      jsonencode(module.support_assign_pending_chat_session),
      jsonencode(module.support_chat_messages_get),
      jsonencode(module.support_chat_sessions_get),
      jsonencode(module.support_chat_sessions_update_open),
      jsonencode(module.support_update_chat_notes),

      jsonencode(module.users_user_list),
      jsonencode(module.users_revoke_access),
      jsonencode(module.users_get_organizations),
      jsonencode(module.users_create_organizations),
      jsonencode(module.users_revoke_access),
      jsonencode(module.users_update_user),

      jsonencode(module.patient_list_patients),
      jsonencode(module.patient_list_appointments),
      jsonencode(module.patient_list_agency_providers),
      jsonencode(module.patient_create_agency_providers),
      jsonencode(module.patient_edit_agency_providers),
      jsonencode(module.patient_create_patients),
      jsonencode(module.patient_patch_patients),
      jsonencode(module.patient_create_appointments),
      jsonencode(module.patient_edit_appointments),
      jsonencode(module.patient_delete_appointments)

    ]))
  }
}

module "route53_configuration" {
  source         = "./modules/certificate-generator"
  api_id         = aws_apigatewayv2_api.gateway.id
  domain_name    = var.domain_name
  domain_prefix  = "api"
  environment    = var.environment
  product_domain = "Sonar Gateway HTTP"
  set_identifier = "api-sonar-gateway-${var.environment}"
  stage_id       = aws_apigatewayv2_stage.v1.id
  hosted_zone_id = var.hosted_zone_id
}
