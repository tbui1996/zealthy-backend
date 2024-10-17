#### users_send start

resource "aws_iam_role" "sonar_service_users_send_role" {
  name               = "sonar_service_users_send_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_users_send_policy" {
  name   = "sonar_service_users_users_send_policy"
  role   = aws_iam_role.sonar_service_users_send_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_users_send_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_users_send_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_send.function_name}:*"]
  }

  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:SendMessage"
    ]
    resources = [var.send_queue_arn]
  }
}
#### users_send end

#### users_receive start
resource "aws_iam_role" "sonar_service_users_receive_role" {
  name               = "sonar_service_users_receive_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_users_receive_policy" {
  name   = "sonar_service_users_users_receive_policy"
  role   = aws_iam_role.sonar_service_users_receive_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_users_receive_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_users_receive_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_receive.function_name}:*"]
  }

  statement {
    actions = [
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl",
      "sqs:ReceiveMessage"
    ]
    resources = [var.receive_queue_arn]
  }
}
#### users_receive end


#### Pre Sign-Up Lambda start
resource "aws_iam_role" "sonar_service_users_pre_sign_up_role" {
  name               = "sonar_service_users_pre_sign_up_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_pre_sign_up_policy" {
  name   = "sonar_service_users_pre_sign_up_policy"
  role   = aws_iam_role.sonar_service_users_pre_sign_up_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_pre_sign_up_policy_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_users_pre_sign_up_vpc" {
  role       = aws_iam_role.sonar_service_users_pre_sign_up_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_users_pre_sign_up_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_pre_sign_up.function_name}:*"]
  }
}
#### Pre Sign-Up Lambda end


#### Post Confirmation Lambda start
resource "aws_iam_role" "sonar_service_users_post_confirmation_role" {
  name               = "sonar_service_users_post_confirmation_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_post_confirmation_policy" {
  name   = "sonar_service_users_post_confirmation_policy"
  role   = aws_iam_role.sonar_service_users_post_confirmation_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_post_confirmation_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_post_confirmation_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_post_confirmation.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminAddUserToGroup"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }
}
#### Post Confirmation Lambda end



#### olive_authorizer_lambda start
resource "aws_iam_role" "sonar_service_users_olive_authorizer_role" {
  name               = "sonar_service_users_olive_authorizer_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_olive_authorizer_policy" {
  name   = "sonar_service_users_olive_authorizer_policy"
  role   = aws_iam_role.sonar_service_users_olive_authorizer_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_olive_authorizer_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_olive_authorizer_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_olive_authorizer.function_name}:*"]
  }

}
#### olive_authorizer_lambda end

#### external_sign_in lambda start
resource "aws_iam_role" "sonar_service_users_external_sign_in_role" {
  name               = "sonar_service_users_external_sign_in_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_external_sign_in_policy" {
  name   = "sonar_service_users_external_sign_in_policy"
  role   = aws_iam_role.sonar_service_users_external_sign_in_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_external_sign_in_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_external_sign_in_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_external_sign_in.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminGetUser",
      "cognito-idp:AdminInitiateAuth"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }
}
#### external_sign_in lambda end

#### external_sign_up lambda start
resource "aws_iam_role" "sonar_service_users_external_sign_up_role" {
  name               = "sonar_service_users_external_sign_up_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_external_sign_up_policy" {
  name   = "sonar_service_users_external_sign_up_policy"
  role   = aws_iam_role.sonar_service_users_external_sign_up_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_external_sign_up_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_external_sign_up_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_external_sign_up.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:SignUp"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }

  statement {
    actions = [
      "dynamodb:GetItem"
    ]
    resources = [aws_dynamodb_table.domain_whitelist.arn]
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
    resources = [aws_kms_key.domain_whitelist_kms_key.arn]
  }
}
#### external_sign_up lambda end

#### external_refresh lambda start
resource "aws_iam_role" "sonar_service_users_external_refresh_role" {
  name               = "sonar_service_users_external_refresh_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_external_refresh_policy" {
  name   = "sonar_service_users_external_refresh_policy"
  role   = aws_iam_role.sonar_service_users_external_refresh_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_external_refresh_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_external_refresh_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_external_refresh.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:InitiateAuth"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }
}
#### external_refresh lambda end

#### define_auth_challenge lambda start
resource "aws_iam_role" "sonar_service_users_define_auth_challenge_role" {
  name               = "sonar_service_users_define_auth_challenge_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_define_auth_challenge_policy" {
  name   = "sonar_service_users_define_auth_challenge_policy"
  role   = aws_iam_role.sonar_service_users_define_auth_challenge_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_define_auth_challenge_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_define_auth_challenge_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_define_auth_challenge.function_name}:*"]
  }
}
#### define_auth_challenge lambda end

#### internal_post_authentication lambda start
resource "aws_iam_role" "sonar_service_users_internal_post_authentication_role" {
  name               = "sonar_service_users_internal_post_authentication_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_internal_post_authentication_policy" {
  name   = "sonar_service_users_internal_post_authentication_policy"
  role   = aws_iam_role.sonar_service_users_internal_post_authentication_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_internal_post_authentication_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_internal_post_authentication_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_internal_post_authentication.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminAddUserToGroup",
      "cognito-idp:AdminListGroupsForUser",
      "cognito-idp:AdminRemoveUserFromGroup",
    ]
    resources = [aws_cognito_user_pool.internals.arn]
  }
}
#### internal_post_authentication lambda end

#### internal_post_confirmation lambda start
resource "aws_iam_role" "sonar_service_users_internal_post_confirmation_role" {
  name               = "sonar_service_users_internal_post_confirmation_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_internal_post_confirmation_policy" {
  name   = "sonar_service_users_internal_post_confirmation_policy"
  role   = aws_iam_role.sonar_service_users_internal_post_confirmation_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_internal_post_confirmation_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_internal_post_confirmation_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_internal_post_confirmation.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminAddUserToGroup",
      "cognito-idp:AdminListGroupsForUser",
      "cognito-idp:AdminRemoveUserFromGroup",
    ]
    resources = [aws_cognito_user_pool.internals.arn]
  }
}
#### internal_post_confirmation lambda end

#### users_list lambda start
resource "aws_iam_role" "sonar_service_users_user_list_role" {
  name               = "sonar_service_users_user_list_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_user_list_policy" {
  name   = "sonar_service_users_user_list_policy"
  role   = aws_iam_role.sonar_service_users_user_list_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_user_list_policy_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_users_user_list_vpc" {
  role       = aws_iam_role.sonar_service_users_user_list_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_users_user_list_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_user_list.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:ListUsers",
      "cognito-idp:AdminListGroupsForUser"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }
}
#### users_list lambda end

#### revoke_access lambda start
resource "aws_iam_role" "sonar_service_users_revoke_access_role" {
  name               = "sonar_service_users_revoke_access_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_revoke_access_policy" {
  name   = "sonar_service_users_revoke_access_policy"
  role   = aws_iam_role.sonar_service_users_revoke_access_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_revoke_access_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_revoke_access_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_revoke_access.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminDisableUser",
      "cognito-idp:AdminGetUser",
      "cognito-idp:AdminListGroupsForUser",
      "cognito-idp:AdminRemoveUserFromGroup",
      "cognito-idp:AdminUserGlobalSignOut"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }

  statement {
    actions = [
      "dynamodb:Query"
    ]
    resources = [var.router_dynamo_arns.websocket_connections]
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
    resources = [var.router_dynamo_arns.websocket_connections_kms_key]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.loop_websocket_api_id}/*"]
  }
}
#### group_assign lambda end

#### get_organizations lambda start
resource "aws_iam_role" "sonar_service_users_get_organizations_role" {
  name               = "sonar_service_users_get_organizations_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_get_organizations_policy" {
  name   = "sonar_service_users_get_organizations_policy"
  role   = aws_iam_role.sonar_service_users_get_organizations_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_get_organizations_policy_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_users_get_organizations_vpc" {
  role       = aws_iam_role.sonar_service_users_get_organizations_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_users_get_organizations_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_get_organizations.function_name}:*"]
  }
}
#### get_organizations lambda end

#### task_get_external_user lambda start
resource "aws_iam_role" "sonar_service_users_task_get_external_user_role" {
  name               = "sonar_service_users_task_get_external_user_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_task_get_external_user_policy" {
  name   = "sonar_service_users_task_get_external_user_policy"
  role   = aws_iam_role.sonar_service_users_task_get_external_user_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_task_get_external_user_policy_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_users_task_get_external_user_vpc" {
  role       = aws_iam_role.sonar_service_users_task_get_external_user_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}


data "aws_iam_policy_document" "sonar_service_users_task_get_external_user_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_task_get_external_user.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminGetUser",
      "cognito-idp:AdminListGroupsForUser"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }
}
#### task_get_external_user lambda end

#### task_get_internal_user lambda start
resource "aws_iam_role" "sonar_service_users_task_get_internal_user_role" {
  name               = "sonar_service_users_task_get_internal_user_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_task_get_internal_user_policy" {
  name   = "sonar_service_users_task_get_internal_user_policy"
  role   = aws_iam_role.sonar_service_users_task_get_internal_user_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_task_get_internal_user_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_task_get_internal_user_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_task_get_internal_user.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminGetUser",
    ]
    resources = [aws_cognito_user_pool.internals.arn]
  }
}
#### task_get_internal_user lambda end

#### update_user lambda start
resource "aws_iam_role" "sonar_service_update_user_role" {
  name               = "sonar_service_update_users_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_update_user_policy" {
  name   = "sonar_service_update_user_policy"
  role   = aws_iam_role.sonar_service_update_user_role.id
  policy = data.aws_iam_policy_document.sonar_service_update_user_policy_document.json
}

#### create_organizations lambda start
resource "aws_iam_role" "sonar_service_users_create_organizations_role" {
  name               = "sonar_service_users_create_organizations_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_create_organizations_policy" {
  name   = "sonar_service_users_create_organizations_policy"
  role   = aws_iam_role.sonar_service_users_create_organizations_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_create_organizations_policy_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_users_create_organizations_vpc" {
  role       = aws_iam_role.sonar_service_users_create_organizations_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_users_create_organizations_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_create_organizations.function_name}:*"]
  }
}
#### create_organizations lambda end

resource "aws_iam_role_policy_attachment" "sonar_service_update_user_vpc" {
  role       = aws_iam_role.sonar_service_update_user_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_update_user_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_get_organizations.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:AdminGetUser",
      "cognito-idp:AdminListGroupsForUser",
      "cognito-idp:AdminUpdateUserAttributes",
      "cognito-idp:AdminAddUserToGroup",
      "cognito-idp:AdminDisableUser",
      "cognito-idp:AdminEnableUser",
      "cognito-idp:AdminRemoveUserFromGroup",
      "cognito-idp:GetGroup"
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }

  statement {
    actions   = ["dynamodb:Query"]
    resources = ["${var.router_dynamo_arns.unconfirmed_websocket_connections}"]
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
    resources = ["${var.router_dynamo_arns.unconfirmed_websocket_connections_kms_key}"]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.loop_unconfirmed_websocket_api_id}/*"]
  }
}
#### update_user lambda end
