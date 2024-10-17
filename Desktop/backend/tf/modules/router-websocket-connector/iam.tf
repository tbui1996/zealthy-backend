#### ${var.name}_forward_role start
resource "aws_iam_role" "route_service_forward_role" {
  name               = "${var.name}_forward_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "route_service_forward_policy" {
  name   = "${var.name}_forward_policy"
  role   = aws_iam_role.route_service_forward_role.id
  policy = data.aws_iam_policy_document.route_service_forward_policy_document.json
}

data "aws_iam_policy_document" "route_service_forward_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.route_service_forward.function_name}:*"]
  }

  statement {
    actions   = ["dynamodb:Scan"]
    resources = [var.router_dynamo_arns.websocket_connections]
  }

  statement {
    actions   = ["dynamodb:PutItem"]
    resources = [var.router_dynamo_arns.pending_messages]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.loop_websocket_api_id}/*"]
  }

  # DeleteMessage and GetQueueURL are used in the service forward lambda, but the others are needed for SQS triggers
  # https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-lambda-function-trigger.html
  statement {
    actions = [
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl",
      "sqs:ReceiveMessage"
    ]
    resources = [aws_sqs_queue.service_send.arn]
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
      var.router_dynamo_arns.websocket_connections_kms_key,
      var.router_dynamo_arns.pending_messages_kms_key,
      aws_kms_key.service_send_kms_key.arn
    ]
  }
}
#### ${var.name}_forward_role end

#### ${var.name}_receive_role start
resource "aws_iam_role" "route_service_receive_role" {
  name               = "${var.name}_receive_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "route_service_receive_policy" {
  name   = "${var.name}_receive_policy"
  role   = aws_iam_role.route_service_receive_role.id
  policy = data.aws_iam_policy_document.route_service_receive_policy_document.json
}

data "aws_iam_policy_document" "route_service_receive_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.lambda.function_name}:*"]
  }

  # GetQueueUrl and SendMessage are used in the service forward lambda, but the others are needed for SQS triggers
  # https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-lambda-function-trigger.html
  statement {
    actions = [
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl",
      "sqs:SendMessage",
      "sqs:ReceiveMessage"
    ]
    resources = [aws_sqs_queue.service_receive.arn]
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
    resources = [aws_kms_key.service_receive_kms_key.arn]
  }
}
#### ${var.name}_receive_role end

#### apigw invocation start
resource "aws_iam_role" "loop_websocket_api_invocation_role" {
  name               = "sonar_loop_websocket_api_${var.name}_invocation_role"
  assume_role_policy = data.aws_iam_policy_document.apigw_execution_document.json
}

resource "aws_iam_role_policy" "loop_websocket_api_invocation_policy" {
  name   = "sonar_loop_websocket_api_${var.name}_invocation_policy"
  role   = aws_iam_role.loop_websocket_api_invocation_role.id
  policy = data.aws_iam_policy_document.loop_websocket_api_invocation_policy_document.json
}

data "aws_iam_policy_document" "loop_websocket_api_invocation_policy_document" {
  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [aws_lambda_function.lambda.arn]
  }
}
#### apigw invocation start
