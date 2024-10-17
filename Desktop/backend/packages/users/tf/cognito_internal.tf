resource "aws_cognito_user_pool" "internals" {
  name = "internals"

  # Custom attribute to be populated with user group from idP
  schema {
    name                = "groups"
    attribute_data_type = "String"
    mutable             = true

    string_attribute_constraints {
      min_length = 1
      max_length = 2048
    }
  }

  lambda_config {
    post_authentication = aws_lambda_function.sonar_service_users_internal_post_authentication.arn
    post_confirmation   = aws_lambda_function.sonar_service_users_internal_post_confirmation.arn
  }

  admin_create_user_config {
    allow_admin_create_user_only = true
  }
}

resource "aws_cognito_user_pool_domain" "internals_web_app" {
  domain       = "sonar-${var.environment}-internals-web-app"
  user_pool_id = aws_cognito_user_pool.internals.id
}

resource "aws_cognito_user_pool_client" "internals_web_app" {
  depends_on = [
    aws_cognito_identity_provider.okta_provider
  ]

  name         = "internals_web_app"
  user_pool_id = aws_cognito_user_pool.internals.id

  allowed_oauth_flows_user_pool_client = true

  supported_identity_providers = ["Okta"]

  allowed_oauth_flows  = ["code"]
  allowed_oauth_scopes = ["email", "openid", "profile"]

  # Where the user is redirected after a successful sign in
  callback_urls = var.callback_urls
  # Where the user is redirected after signing out
  logout_urls = var.logout_urls
}

resource "okta_app_oauth" "internals_okta_app_integration" {
  label             = "sonar-${var.environment}"
  auto_key_rotation = true
  grant_types       = ["authorization_code"]
  redirect_uris     = ["https://${aws_cognito_user_pool_domain.internals_web_app.domain}.auth.${data.aws_region.current.name}.amazoncognito.com/oauth2/idpresponse"]
  response_types    = ["code"]
  type              = "web"
  lifecycle {
    ignore_changes = [users, groups]
  }
  groups_claim {
    type        = "FILTER"
    filter_type = "STARTS_WITH"
    name        = "groups"
    value       = "internals"
  }
}

resource "okta_app_bookmark" "internals_okta_app_bookmark" {
  label = var.environment == "prod" ? "Sonar" : "Sonar (${var.environment})"
  url   = "https://${aws_cognito_user_pool_domain.internals_web_app.domain}.auth.${data.aws_region.current.name}.amazoncognito.com/login?response_type=code&client_id=${aws_cognito_user_pool_client.internals_web_app.id}&scope=${join("+", aws_cognito_user_pool_client.internals_web_app.allowed_oauth_scopes)}&redirect_uri=${one(aws_cognito_user_pool_client.internals_web_app.callback_urls)}"
  logo  = "${path.module}/assets/03-Sonar-LT-128.png"
  lifecycle {
    ignore_changes = [users, groups]
  }
}

resource "okta_group" "groups" {
  for_each    = var.internal_group_names
  name        = "internals_${each.key}.${var.environment}"
  description = "Group for ${each.key} in ${var.environment}"
  lifecycle {
    ignore_changes = [users]
  }
}

// TODO: does everyone who ever uses web get general_support? Or should there also be guest role
resource "okta_user_group_memberships" "assign_group_to_dev" {
  # developer accounts only!!!
  count   = var.live_env ? 0 : 1
  user_id = var.okta_user_id
  # when we support multiple groups, e.g. Admin, we'll likely need to update this to the super user
  # assumption  that each sonar user will belong to only one group
  groups = [
    okta_group.groups["program_manager"].id
  ]
}

resource "okta_app_group_assignment" "sonar_integration" {
  for_each = okta_group.groups
  app_id   = okta_app_oauth.internals_okta_app_integration.id
  group_id = each.value.id
}

resource "okta_app_group_assignment" "sonar_integration_bookmark" {
  for_each = okta_group.groups
  app_id   = okta_app_bookmark.internals_okta_app_bookmark.id
  group_id = each.value.id
}

resource "aws_cognito_identity_provider" "okta_provider" {
  user_pool_id  = aws_cognito_user_pool.internals.id
  provider_name = "Okta"
  provider_type = "OIDC"

  provider_details = {
    client_id                 = okta_app_oauth.internals_okta_app_integration.client_id
    client_secret             = okta_app_oauth.internals_okta_app_integration.client_secret
    attributes_request_method = "GET"
    oidc_issuer               = "https://${var.okta_org_name}.${var.okta_base_url}"
    authorize_scopes          = "openid email profile groups"
  }
  attribute_mapping = {
    # openid scope
    username = "sub"

    # email scope
    email          = "email"
    email_verified = "email_verified"

    # profile scope
    name = "name"

    # groups scope
    "custom:groups" = "groups"
  }
}

resource "aws_lambda_permission" "cognito_internal_post_authentication" {
  statement_id  = "AllowExecutionFromUserPool"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sonar_service_users_internal_post_authentication.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.internals.arn
}

resource "aws_lambda_permission" "cognito_internal_post_confirmation" {
  statement_id  = "AllowExecutionFromUserPool"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sonar_service_users_internal_post_confirmation.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.internals.arn
}
