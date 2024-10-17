#### route_connect start
resource "aws_iam_role" "global_route_connect_role" {
  name               = "global_route_connect_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "global_route_connect_policy" {
  name = "global_route_connect_policy"
  role = aws_iam_role.global_route_connect_role.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ],
          Resource : "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.internal_websocket_connect.function_name}:*",
          Effect : "Allow"
        },
        {
          Sid : "",
          Effect : "Allow",
          Action : [
            "dynamodb:PutItem"
          ],
          Resource : "${var.router_dynamo_arns.websocket_connections_internal}"
        },
        {
          Sid : "",
          Effect : "Allow",
          Action : [
            "events:PutEvents"
          ],
          Resource : "${var.event_bus_arn}"
        },
        {
          Sid : "",
          Effect : "Allow",
          Action : [
            "kms:DescribeKey",
            "kms:Encrypt",
            "kms:Decrypt",
            "kms:ReEncryptTo",
            "kms:ReEncryptFrom",
            "kms:GenerateDataKey",
            "kms:GenerateDataKeyWithoutPlaintext"
          ],
          Resource : "${var.router_dynamo_arns.websocket_connections_internal_kms_key}"
        }
      ]
  })
}
#### route_connect end

#### route_disconnect start
resource "aws_iam_role" "global_route_disconnect_role" {
  name               = "global_route_disconnect_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "global_route_disconnect_policy" {
  name = "global_route_disconnect_policy"
  role = aws_iam_role.global_route_disconnect_role.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ],
          Resource : "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.internal_websocket_disconnect.function_name}:*",
          Effect : "Allow"
        },
        {
          Sid : "",
          Effect : "Allow",
          Action : [
            "dynamodb:DeleteItem"
          ],
          Resource : "${var.router_dynamo_arns.websocket_connections_internal}"
        },
        {
          Sid : "",
          Effect : "Allow",
          Action : [
            "kms:DescribeKey",
            "kms:Encrypt",
            "kms:Decrypt",
            "kms:ReEncryptTo",
            "kms:ReEncryptFrom",
            "kms:GenerateDataKey",
            "kms:GenerateDataKeyWithoutPlaintext"
          ],
          Resource : "${var.router_dynamo_arns.websocket_connections_internal_kms_key}"
        }
      ]
  })
}
#### route_disconnect end

### global_task_is_internal_user_online start
resource "aws_iam_role" "global_task_is_internal_user_online_role" {
  name               = "global_task_is_internal_user_online_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "global_task_is_internal_user_online_policy" {
  name   = "global_task_is_internal_user_online_policy"
  role   = aws_iam_role.global_task_is_internal_user_online_role.id
  policy = data.aws_iam_policy_document.global_task_is_internal_user_online_policy_document.json
}

data "aws_iam_policy_document" "global_task_is_internal_user_online_policy_document" {
  policy_id = "global_task_is_internal_user_online_policy_document"

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.global_task_is_internal_user_online.function_name}:*"]
  }

  statement {
    actions = [
      "dynamodb:Query",
    ]
    resources = [var.router_dynamo_arns.websocket_connections_internal]
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
### global_task_is_internal_user_online end
