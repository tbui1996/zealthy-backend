data "aws_iam_policy_document" "logging_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/sonar_service_support_*:*",
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/offline_message_notifier_*:*",
    ]
  }
}

data "aws_iam_policy_document" "apigateway_web_invoke_api_policy_document" {
  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/*/support/*"
    ]
  }
}

data "aws_iam_policy_document" "apigateway_loop_invoke_api_policy_document" {
  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.aws_apigatewayv2_loop_gateway_id}/${var.aws_apigatewayv2_loop_gateway_stage_name}/*/support/*"
    ]
  }
}

data "aws_iam_policy_document" "apigateway_websocket_web_invoke_api_policy_document" {
  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.webapp_websocket_api_id}/*"
    ]
  }
}

data "aws_iam_policy_document" "apigateway_websocket_loop_invoke_api_policy_document" {
  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.loop_websocket_api_id}/*"
    ]
  }
}

data "aws_iam_policy_document" "chat_messages_kms_key_policy_document" {
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
      aws_kms_key.messages_kms_key.arn,
    ]
  }
}

## -- Assign Pending Chat Start --
resource "aws_iam_role" "sonar_service_support_assign_pending_chat_session" {
  name               = "sonar_service_support_assign_pending_chat_session"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_assign_pending_chat_session_vpc" {
  role       = aws_iam_role.sonar_service_support_assign_pending_chat_session.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_assign_pending_chat_session_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.chat_messages_kms_key_policy_document.json,
  ]
  statement {
    actions = [
      "dynamodb:Query",
    ]
    resources = [
      aws_dynamodb_table.messages.arn,
      "${aws_dynamodb_table.messages.arn}/*"
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_assign_pending_chat_session_role_policy" {
  name   = "sonar_service_assign_pending_chat_session_policy"
  policy = data.aws_iam_policy_document.sonar_service_assign_pending_chat_session_document.json
  role   = aws_iam_role.sonar_service_support_assign_pending_chat_session.id
}
## -- Assign Pending Chat End --

## -- Chat Messages Get Start --
resource "aws_iam_role" "sonar_service_support_chat_messages_get" {
  name               = "sonar_service_support_chat_messages_get"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_chat_messages_get_vpc" {
  role       = aws_iam_role.sonar_service_support_chat_messages_get.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_chat_messages_get_policy" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
    data.aws_iam_policy_document.chat_messages_kms_key_policy_document.json
  ]
  statement {
    actions = ["dynamodb:Query"]
    resources = [
      aws_dynamodb_table.messages.arn,
      "${aws_dynamodb_table.messages.arn}/*"
    ]
  }
}

resource "aws_iam_role_policy" "support_chat_messages_get_role_policy" {
  name   = "support_chat_messages_get_policy"
  role   = aws_iam_role.sonar_service_support_chat_messages_get.id
  policy = data.aws_iam_policy_document.sonar_service_support_chat_messages_get_policy.json
}
## -- Chat Messages Get End --

## -- Chat Session Create Start --
resource "aws_iam_role" "support_chat_session_create" {
  name               = "support_chat_session_create"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "support_chat_session_create_vpc" {
  role       = aws_iam_role.support_chat_session_create.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "support_chat_session_create_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json
  ]
}

resource "aws_iam_role_policy" "support_chat_session_create_policy" {
  name   = "support_chat_session_create_policy"
  role   = aws_iam_role.support_chat_session_create.id
  policy = data.aws_iam_policy_document.support_chat_session_create_policy_document.json
}
## -- Chat Session Create End --

## -- Chat Session Update Start --
resource "aws_iam_role" "sonar_service_support_chat_session_update" {
  name               = "sonar_service_support_chat_session_update"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_chat_session_update_vpc" {
  role       = aws_iam_role.sonar_service_support_chat_session_update.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_chat_session_update_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
  ]
}

resource "aws_iam_role_policy" "support_chat_session_update_role_policy" {
  name   = "support_chat_session_update_policy"
  role   = aws_iam_role.sonar_service_support_chat_session_update.id
  policy = data.aws_iam_policy_document.sonar_service_support_chat_session_update_policy_document.json
}
## -- Chat Session Update End --

## -- Chat Session Get Start --
resource "aws_iam_role" "sonar_service_support_chat_sessions_get" {
  name               = "sonar_service_support_chat_sessions_get"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_chat_sessions_get_vpc" {
  role       = aws_iam_role.sonar_service_support_chat_sessions_get.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_chat_sessions_get_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
  ]
}

resource "aws_iam_role_policy" "support_chat_sessions_get_role_policy" {
  name   = "support_chat_sessions_get_role_policy"
  role   = aws_iam_role.sonar_service_support_chat_sessions_get.id
  policy = data.aws_iam_policy_document.sonar_service_support_chat_sessions_get_policy_document.json
}
## -- Chat Session Get End --

## -- Pending Chat Session Create Start --
resource "aws_iam_role" "support_pending_chat_session_create" {
  name               = "support_pending_chat_session_create"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "support_pending_chat_session_create_vpc" {
  role       = aws_iam_role.support_pending_chat_session_create.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "support_pending_chat_session_create_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_web_invoke_api_policy_document.json,
  ]
  statement {
    actions = ["dynamodb:Query"]
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
    resources = [var.router_dynamo_arns.websocket_connections_internal_kms_key]
  }
}

resource "aws_iam_role_policy" "support_pending_chat_session" {
  name   = "support_pending_chat_session"
  role   = aws_iam_role.support_pending_chat_session_create.id
  policy = data.aws_iam_policy_document.support_pending_chat_session_create_policy_document.json
}
## -- Pending Chat Session Create End --

## -- Pending Chat Session Get Start --
resource "aws_iam_role" "support_pending_chat_sessions_get" {
  name               = "support_pending_chat_sessions_get"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "support_pending_chat_sessions_get_vpc" {
  role       = aws_iam_role.support_pending_chat_sessions_get.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "support_pending_chat_sessions_get_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
  ]
}

resource "aws_iam_role_policy" "support_pending_chat_sessions_get_dynamodb_policy" {
  name   = "support_pending_chat_sessions_get_policy"
  role   = aws_iam_role.support_pending_chat_sessions_get.id
  policy = data.aws_iam_policy_document.support_pending_chat_sessions_get_policy_document.json
}
## -- Pending Chat Session Get End --

## -- Support Receive Start --
resource "aws_iam_role" "sonar_service_support_receive" {
  name               = "sonar_service_support_receive"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_receive_vpc" {
  role       = aws_iam_role.sonar_service_support_receive.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_receive_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_loop_invoke_api_policy_document.json,
    data.aws_iam_policy_document.chat_messages_kms_key_policy_document.json,
  ]
  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes"
    ]
    resources = [
      var.receive_queue_arn
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
      var.receive_queue_kms_arn
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_support_receive_policy" {
  name   = "sonar_service_support_receive_policy"
  role   = aws_iam_role.sonar_service_support_receive.id
  policy = data.aws_iam_policy_document.sonar_service_support_receive_policy_document.json
}
## -- Support Receive End --

## -- Support Send Start --
resource "aws_iam_role" "sonar_service_support_send" {
  name               = "sonar_service_support_send"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_send_vpc" {
  role       = aws_iam_role.sonar_service_support_send.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_send_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_websocket_loop_invoke_api_policy_document.json,
    data.aws_iam_policy_document.chat_messages_kms_key_policy_document.json
  ]
  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:SendMessage"
    ]
    resources = [
      var.send_queue_arn
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
      var.send_queue_kms_arn
    ]
  }
  statement {
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [
      aws_dynamodb_table.messages.arn
    ]
  }
  statement {
    actions = ["events:PutEvents"]
    resources = [
      var.event_bus_arn
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_support_send_policy" {
  name   = "sonar_service_support_send_policy"
  policy = data.aws_iam_policy_document.sonar_service_support_send_policy_document.json
  role   = aws_iam_role.sonar_service_support_send.id
}
## -- Support Send End --

## -- Support update chat notes start --
resource "aws_iam_role" "sonar_service_support_update_chat_notes" {
  name               = "sonar_service_support_update_chat_notes"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_update_chat_notes_vpc" {
  role       = aws_iam_role.sonar_service_support_update_chat_notes.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_update_chat_notes_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
  ]
}

resource "aws_iam_role_policy" "support_update_chat_notes_role_policy" {
  name   = "support_update_chat_notes_role_policy"
  role   = aws_iam_role.sonar_service_support_update_chat_notes.id
  policy = data.aws_iam_policy_document.sonar_service_support_update_chat_notes_policy_document.json
}
## -- Support update chat notes end --

## -- Support submit feedback start --
resource "aws_iam_role" "sonar_service_support_submit_feedback" {
  name               = "sonar_service_support_submit_feedback"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_submit_feedback_vpc" {
  role       = aws_iam_role.sonar_service_support_submit_feedback.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_submit_feedback_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
  ]

  statement {
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [
      aws_dynamodb_table.feedback.arn
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
      aws_kms_key.feedback_kms_key.arn,
    ]
  }
  statement {
    actions = [
      "ses:SendTemplatedEmail",
    ]
    resources = ["arn:aws:ses:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"]
  }
}

resource "aws_iam_role_policy" "support_submit_feedback_role_policy" {
  name   = "support_submit_feedback_role_policy"
  role   = aws_iam_role.sonar_service_support_submit_feedback.id
  policy = data.aws_iam_policy_document.sonar_service_support_submit_feedback_policy_document.json
}
## -- Support submit feedback end --

### Start Consume Event State Handler Permissions ###
resource "aws_iam_role" "offline_message_notifier_task_consume_send_message_event_role" {
  name               = "offline_message_notifier_task_consume_send_message_event_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "offline_message_notifier_task_consume_send_message_event_role_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.offline_message_notifier_task_consume_send_message_event_fn.function_name}:*",
    ]
  }
}

resource "aws_iam_role_policy" "offline_message_notifier_task_consume_send_message_event_policy" {
  name = "offline_message_notifier_task_consume_send_message_event_policy-policy"
  role = aws_iam_role.offline_message_notifier_task_consume_send_message_event_role.id

  policy = data.aws_iam_policy_document.offline_message_notifier_task_consume_send_message_event_role_policy_document.json
}
### End Consume Event State Handler Permissions ###

### Start Reset Offline Message Permissions ###
resource "aws_iam_role" "offline_message_notifier_task_reset_offline_message_role" {
  name               = "offline_message_notifier_task_reset_offline_message_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "offline_message_notifier_task_reset_offline_message_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
  ]

  statement {
    actions = [
      "dynamodb:DeleteItem",
    ]
    resources = [
      aws_dynamodb_table.offline_messages.arn
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
      aws_kms_key.offline_messages_kms_key.arn,
    ]
  }
}

resource "aws_iam_role_policy" "offline_message_task_reset_offline_message_policy" {
  name = "offline_message_task_reset_offline_message_policy"
  role = aws_iam_role.offline_message_notifier_task_reset_offline_message_role.id

  policy = data.aws_iam_policy_document.offline_message_notifier_task_reset_offline_message_document.json
}
### End Reset Offline Message Permissions ###

### Start Record Offline Message Permissions ###
resource "aws_iam_role" "offline_message_notifier_task_record_offline_message_role" {
  name               = "offline_message_notifier_task_record_offline_message_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "offline_message_notifier_task_record_offline_message_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
  ]

  statement {
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [
      aws_dynamodb_table.offline_messages.arn
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
      aws_kms_key.offline_messages_kms_key.arn,
    ]
  }
}

resource "aws_iam_role_policy" "offline_message_task_record_offline_message_policy" {
  name = "offline_message_task_record_offline_message_policy"
  role = aws_iam_role.offline_message_notifier_task_record_offline_message_role.id

  policy = data.aws_iam_policy_document.offline_message_notifier_task_record_offline_message_document.json
}
### End Record Offline Message Permissions ###

### Start Trigger Email Permissions ###
resource "aws_iam_role" "offline_message_notifier_task_trigger_offline_email_role" {
  name               = "offline_message_notifier_task_trigger_offline_email_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "offline_message_notifier_task_trigger_offline_email_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
  ]

  statement {
    actions = [
      "dynamodb:UpdateItem",
    ]
    resources = [
      aws_dynamodb_table.offline_messages.arn
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
      aws_kms_key.offline_messages_kms_key.arn,
    ]
  }

  statement {
    actions = [
      "ses:SendEmail",
      "ses:SendTemplatedEmail",
    ]
    resources = ["arn:aws:ses:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"]
  }
}

resource "aws_iam_role_policy" "offline_message_task_trigger_offline_email_policy" {
  name = "offline_message_task_trigger_offline_email_policy"
  role = aws_iam_role.offline_message_notifier_task_trigger_offline_email_role.id

  policy = data.aws_iam_policy_document.offline_message_notifier_task_trigger_offline_email_document.json
}
### End Trigger Email Permissions ###

## -- Support chat session star start --
resource "aws_iam_role" "sonar_service_support_chat_session_star" {
  name               = "sonar_service_support_chat_session_star"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_chat_session_star_vpc" {
  role       = aws_iam_role.sonar_service_support_chat_session_star.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_chat_session_star_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json,
  ]
}

resource "aws_iam_role_policy" "support_chat_session_star_role_policy" {
  name   = "support_chat_session_star_role_policy"
  role   = aws_iam_role.sonar_service_support_chat_session_star.id
  policy = data.aws_iam_policy_document.sonar_service_support_chat_session_star_policy_document.json
}
## -- Support chat session star end --

resource "aws_iam_role" "sonar_service_support_loop_online_internal_users" {
  name               = "sonar_service_support_loop_online_internal_users"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_support_loop_online_internal_users" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
  ]
  statement {
    actions = [
      "dynamodb:Scan",
    ]
    resources = [
      var.router_dynamo_arns.websocket_connections_internal
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
}

resource "aws_iam_role_policy" "sonar_service_support_loop_online_internal_users" {
  name   = "sonar_service_support_loop_online_internal_users"
  role   = aws_iam_role.sonar_service_support_loop_online_internal_users.id
  policy = data.aws_iam_policy_document.sonar_service_support_loop_online_internal_users.json
}