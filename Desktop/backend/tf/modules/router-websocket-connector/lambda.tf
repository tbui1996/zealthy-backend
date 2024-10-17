locals {
  variables = {
    "WEBSOCKET_URL" = "ws-sonar.${var.domain_name}"
    "CONTEXT"       = var.name
    "RECEIVE_QUEUE" = aws_sqs_queue.service_receive.name
    "SEND_QUEUE"    = aws_sqs_queue.service_send.name
  }
}

resource "aws_lambda_function" "route_service_forward" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_router_${var.name}_forward"
  role          = aws_iam_role.route_service_forward_role.arn
  timeout       = "29"

  image_config {
    entry_point = ["/lambda/service_forward"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = local.variables
  }
}

resource "aws_apigatewayv2_route" "apigw_route" {
  api_id    = var.loop_websocket_api_id
  route_key = var.name
  target    = "integrations/${aws_apigatewayv2_integration.apigw_integration.id}"
}

resource "aws_lambda_function" "lambda" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_router_${var.name}_receive"
  role          = aws_iam_role.route_service_receive_role.arn
  timeout       = 29

  image_config {
    entry_point = ["/lambda/service_receive"]
  }

  environment {
    variables = local.variables
  }

  tracing_config {
    mode = "Active"
  }
}

# Lambda Permissions
resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.loop_websocket_api_execution_arn}/*/*"
}

# API Gateway Integration
resource "aws_apigatewayv2_integration" "apigw_integration" {
  api_id                    = var.loop_websocket_api_id
  integration_type          = "AWS_PROXY"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.lambda.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
  credentials_arn           = aws_iam_role.loop_websocket_api_invocation_role.arn
  content_handling_strategy = "CONVERT_TO_TEXT"
}

