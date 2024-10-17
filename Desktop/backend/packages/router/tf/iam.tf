#### api_broadcast_lambda start
resource "aws_iam_role" "api_broadcast_lambda_role" {
  name               = "api_broadcast_lambda_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "api_broadcast_lambda_policy" {
  name   = "api_broadcast_lambda_policy"
  role   = aws_iam_role.api_broadcast_lambda_role.id
  policy = data.aws_iam_policy_document.api_broadcast_lambda_policy_document.json
}

# TODO: revisit broadcast permissions when reworked in SONAR-325. These will probably be able to be tightened up.
data "aws_iam_policy_document" "api_broadcast_lambda_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.api_broadcast_lambda.function_name}:*"]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.loop_websocket_api_id}/*"]
  }

  statement {
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]
    resources = [aws_dynamodb_table.websocket_connections.arn]
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
    resources = [aws_kms_key.websocket_connections_kms_key.arn]
  }

  statement {
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]
    resources = [aws_dynamodb_table.pending_messages.arn]
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
    resources = [aws_kms_key.pending_messages_kms_key.arn]
  }
}
#### api_broadcast_lambda end

# DEPRECATED: remove as part of SONAR-325
#### api_users_list start
resource "aws_iam_role" "api_users_list_role" {
  name               = "sonar_service_${var.name}_api_users_list_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "api_users_list_policy" {
  name   = "logging_policy"
  role   = aws_iam_role.api_users_list_role.id
  policy = data.aws_iam_policy_document.api_users_list_policy_document.json
}

data "aws_iam_policy_document" "api_users_list_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.api_users_list.function_name}:*"]
  }
}

resource "aws_iam_role_policy_attachment" "AWSLambdaVPCAccessExecutionRole" {
  role       = aws_iam_role.api_users_list_role.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
#### api_users_list end

#### route_connect start
resource "aws_iam_role" "router_route_connect_role" {
  name               = "sonar_service_${var.name}_router_route_connect_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "router_route_connect_policy" {
  name   = "logging_policy"
  role   = aws_iam_role.router_route_connect_role.id
  policy = data.aws_iam_policy_document.router_route_connect_policy_document.json
}

data "aws_iam_policy_document" "router_route_connect_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_router_connect.function_name}:*"]
  }

  statement {
    actions   = ["dynamodb:PutItem"]
    resources = [aws_dynamodb_table.websocket_connections.arn]
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
    resources = [aws_kms_key.websocket_connections_kms_key.arn]
  }

  statement {
    actions = ["events:PutEvents"]
    resources = [
      var.event_bus_arn
    ]
  }
}
#### route_connect end

#### route_disconnect start
resource "aws_iam_role" "router_route_disconnect_role" {
  name               = "router_route_disconnect_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "router_route_disconnect_policy" {
  name   = "logging_policy"
  role   = aws_iam_role.router_route_disconnect_role.id
  policy = data.aws_iam_policy_document.router_route_disconnect_policy_document.json
}

data "aws_iam_policy_document" "router_route_disconnect_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_router_disconnect.function_name}:*"]
  }

  statement {
    actions   = ["dynamodb:DeleteItem"]
    resources = [aws_dynamodb_table.websocket_connections.arn]
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
    resources = [aws_kms_key.websocket_connections_kms_key.arn]
  }
}
#### route_disconnect end

#### route_message start
resource "aws_iam_role" "router_route_message_role" {
  name               = "router_route_message_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "router_route_message_policy" {
  name   = "logging_policy"
  role   = aws_iam_role.router_route_message_role.id
  policy = data.aws_iam_policy_document.router_route_message_policy_document.json
}

data "aws_iam_policy_document" "router_route_message_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_router_message.function_name}:*"]
  }

  statement {
    actions   = ["dynamodb:Scan"]
    resources = [aws_dynamodb_table.websocket_connections.arn]
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
    resources = [aws_kms_key.websocket_connections_kms_key.arn]
  }
}
#### route_message end

#### sonar_service_router_receive start
resource "aws_iam_role" "sonar_service_router_receive_role" {
  name               = "sonar_service_router_receive_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_router_receive_policy" {
  name   = "sonar_service_router_receive_policy"
  role   = aws_iam_role.sonar_service_router_receive_role.id
  policy = data.aws_iam_policy_document.sonar_service_router_receive_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_router_receive_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_router_receive.function_name}:*"]
  }

  statement {
    actions = [
      "dynamodb:Scan"
    ]
    resources = [aws_dynamodb_table.websocket_connections.arn]
  }

  statement {
    actions = [
      "dynamodb:Query"
    ]
    resources = [aws_dynamodb_table.pending_messages.arn]
  }

  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:DeleteMessage",
      "sqs:ReceiveMessage",
      "sqs:GetQueueAttributes"
    ]
    resources = [var.receive_queue_arn]
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
      var.receive_queue_kms_arn,
      aws_kms_key.websocket_connections_kms_key.arn,
      aws_kms_key.pending_messages_kms_key.arn
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.loop_websocket_api_id}/*"]
  }
}
#### sonar_service_router_receive end

#### unconfirmed_route_connect start
resource "aws_iam_role" "router_unconfirmed_route_connect_role" {
  name               = "router_unconfirmed_route_connect_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "router_unconfirmed_route_connect_policy" {
  name   = "router_unconfirmed_route_connect_policy"
  role   = aws_iam_role.router_unconfirmed_route_connect_role.id
  policy = data.aws_iam_policy_document.router_unconfirmed_route_connect_policy_document.json
}

data "aws_iam_policy_document" "router_unconfirmed_route_connect_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_router_unconfirmed_connect.function_name}:*"]
  }

  statement {
    actions   = ["dynamodb:PutItem"]
    resources = [aws_dynamodb_table.unconfirmed_websocket_connections.arn]
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
    resources = [aws_kms_key.unconfirmed_websocket_connections_kms_key.arn]
  }
}
#### unconfirmed_route_connect end

#### unconfirmed_route_disconnect start
resource "aws_iam_role" "router_unconfirmed_route_disconnect_role" {
  name               = "router_unconfirmed_route_disconnect_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "router_unconfirmed_route_disconnect_policy" {
  name   = "router_unconfirmed_route_disconnect_policy"
  role   = aws_iam_role.router_unconfirmed_route_disconnect_role.id
  policy = data.aws_iam_policy_document.router_unconfirmed_route_disconnect_policy_document.json
}

data "aws_iam_policy_document" "router_unconfirmed_route_disconnect_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_unconfirmed_disconnect.function_name}:*"]
  }

  statement {
    actions   = ["dynamodb:DeleteItem"]
    resources = [aws_dynamodb_table.unconfirmed_websocket_connections.arn]
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
    resources = [aws_kms_key.unconfirmed_websocket_connections_kms_key.arn]
  }
}
#### unconfirmed_route_disconnect end

### sonar_service_task_is_external_user_online start
resource "aws_iam_role" "router_task_is_external_user_online_role" {
  name               = "router_task_is_external_user_online_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "router_task_is_external_user_online_policy" {
  name   = "router_task_is_external_user_online_policy"
  role   = aws_iam_role.router_task_is_external_user_online_role.id
  policy = data.aws_iam_policy_document.router_task_is_external_user_online_policy_document.json
}

data "aws_iam_policy_document" "router_task_is_external_user_online_policy_document" {
  policy_id = "router_task_is_external_user_online_policy_document"

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_router_task_is_external_user_online.function_name}:*"]
  }

  statement {
    actions = [
      "dynamodb:Query",
    ]
    resources = [aws_dynamodb_table.websocket_connections.arn]
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
    resources = [aws_kms_key.websocket_connections_kms_key.arn]
  }
}
### sonar_service_task_is_external_user_online end