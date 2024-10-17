resource "aws_cognito_user_pool" "externals" {
  name = "externals"

  # Whether email addresses or phone numbers can be specified as usernames when a user signs up
  # This allows users to sign up with email as the username
  username_attributes = ["email"]

  username_configuration {
    case_sensitive = false
  }

  # We don't need users to separately confirm their email address
  auto_verified_attributes = ["email"]

  admin_create_user_config {
    allow_admin_create_user_only = false
  }

  lambda_config {
    pre_sign_up           = aws_lambda_function.sonar_service_users_pre_sign_up.arn
    post_confirmation     = aws_lambda_function.sonar_service_users_post_confirmation.arn
    define_auth_challenge = aws_lambda_function.sonar_service_users_define_auth_challenge.arn
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "firstName"
    required                 = false

    string_attribute_constraints {
      max_length = "256"
      min_length = "1"
    }
  }
  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "lastName"
    required                 = false

    string_attribute_constraints {
      max_length = "256"
      min_length = "1"
    }
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "organization"
    required                 = false

    string_attribute_constraints {
      max_length = "256"
      min_length = "1"
    }
  }
}

resource "aws_cognito_user_pool_client" "externals" {
  name = "externals"

  user_pool_id    = aws_cognito_user_pool.externals.id
  generate_secret = false

  explicit_auth_flows = [
    "CUSTOM_AUTH_FLOW_ONLY"
  ]
}

resource "aws_lambda_permission" "cognito_pre_sign_up" {
  statement_id  = "AllowExecutionFromUserPool"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sonar_service_users_pre_sign_up.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.externals.arn
}

resource "aws_lambda_permission" "cognito_post_confirmation" {
  statement_id  = "AllowExecutionFromUserPool"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sonar_service_users_post_confirmation.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.externals.arn
}

resource "aws_lambda_permission" "cognito_define_auth_challenge" {
  statement_id  = "AllowExecutionFromUserPool"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sonar_service_users_define_auth_challenge.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.externals.arn
}
