resource "aws_iam_role" "sonar_service_support_loop_send" {
  name               = "sonar_service_support_loop_send"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_loop_send_vpc" {
  role       = aws_iam_role.sonar_service_support_loop_send.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_loop_send_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_loop_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
    data.aws_iam_policy_document.chat_messages_kms_key_policy_document.json,
  ]
  statement {
    actions = ["dynamodb:PutItem"]
    resources = [
      aws_dynamodb_table.messages.arn
    ]
  }
  statement {
    actions = [
      "dynamodb:Query",
    ]
    resources = [
      var.router_dynamo_arns.websocket_connections_internal,
      "${var.router_dynamo_arns.websocket_connections_internal}/index/UserGroupIndex"
    ]
  }
  statement {
    actions = [
      "kms:DescribeKey",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncryptTo",
      "kms:ReEncryptFrom",
      "kms:GenerateDataKey",
      "kms:GenerateDataKeyWithoutPlaintext"
    ]
    resources = [
      var.router_dynamo_arns.websocket_connections_internal_kms_key
    ]
  }
  statement {
    actions = ["events:PutEvents"]
    resources = [
      var.event_bus_arn
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_support_loop_send" {
  name   = "sonar_service_support_loop_send"
  role   = aws_iam_role.sonar_service_support_loop_send.id
  policy = data.aws_iam_policy_document.sonar_service_loop_send_policy_document.json
}

resource "aws_lambda_function" "sonar_service_support_loop_send" {
  function_name = "sonar_service_support_loop_send"
  role          = aws_iam_role.sonar_service_support_loop_send.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/loop_send"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION    = data.aws_region.current.name
      WEBSOCKET_URL = "https://ws-sonar-internal.${var.domain_name}"
      "HOST"        = var.db_host
      "NAME"        = var.db_name
      "PASSWORD"    = var.db_password
      "PORT"        = var.db_port
      "USER"        = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }

}
