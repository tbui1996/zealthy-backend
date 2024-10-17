# Give CloudWatch Log Group KMS permission
# https://github.com/hashicorp/terraform-provider-aws/issues/8042
# https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/encrypt-log-data-kms.html
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key#policy
locals {
  loop_websocket_api_logs_name               = "sonar_loop_websocket_api_logs"
  loop_websocket_api_logs_kms_key_alias_name = "alias/loop_websocket_api_logs_kms_key"
}

resource "aws_kms_key" "loop_websocket_api_logs_kms_key" {
  description             = "KMS key for sonar_loop_websocket_api_logs CloudWatch log group"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.loop_websocket_api_logs_kms_key_document.json
}

resource "aws_kms_alias" "loop_websocket_api_logs_kms_key" {
  name          = local.loop_websocket_api_logs_kms_key_alias_name
  target_key_id = aws_kms_key.loop_websocket_api_logs_kms_key.id
}

data "aws_iam_policy_document" "loop_websocket_api_logs_kms_key_document" {
  policy_id = "sonar-loop-websocket-api-logs-kms-key"
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
      values   = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:${local.loop_websocket_api_logs_name}"]
    }
  }
}

resource "aws_cloudwatch_log_group" "loop_websocket_api_logs" {
  name       = local.loop_websocket_api_logs_name
  kms_key_id = aws_kms_key.loop_websocket_api_logs_kms_key.arn
}

resource "aws_apigatewayv2_api" "loop_websocket_api" {
  name                       = "sonar_loop_websocket_api"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "loop_websocket_v1" {
  api_id        = aws_apigatewayv2_api.loop_websocket_api.id
  name          = var.environment
  deployment_id = aws_apigatewayv2_deployment.websocket_deploy.id

  default_route_settings {
    throttling_rate_limit  = 1000
    throttling_burst_limit = 1000
  }

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.loop_websocket_api_logs.arn
    # Common log format
    format = "$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] \"$context.httpMethod $context.resourcePath $context.protocol\" $context.status $context.responseLength $context.requestId"
  }
}

# API Gateway Deployment
resource "aws_apigatewayv2_deployment" "websocket_deploy" {
  api_id = aws_apigatewayv2_api.loop_websocket_api.id

  lifecycle {
    create_before_destroy = true
  }

  triggers = {
    redeployment = sha1(join(",", [
      jsonencode(module.router_connect),
      jsonencode(module.router_disconnect),
      jsonencode(module.router_message),
    ]))
  }
}

module "loop_websocket_route53_configuration" {
  source         = "./modules/certificate-generator"
  api_id         = aws_apigatewayv2_api.loop_websocket_api.id
  domain_name    = var.domain_name
  domain_prefix  = "ws-sonar"
  environment    = var.environment
  product_domain = "Sonar Websockets"
  set_identifier = "ws-sonar-${var.environment}"
  stage_id       = aws_apigatewayv2_stage.loop_websocket_v1.id
  hosted_zone_id = var.hosted_zone_id
}
