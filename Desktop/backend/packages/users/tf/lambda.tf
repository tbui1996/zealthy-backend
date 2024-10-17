resource "aws_lambda_function" "sonar_service_users_send" {
  function_name = "sonar_service_users_send"
  role          = aws_iam_role.sonar_service_users_send_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/send"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_receive" {
  function_name = "sonar_service_users_receive"
  role          = aws_iam_role.sonar_service_users_receive_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/receive"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_pre_sign_up" {
  function_name = "sonar_service_users_pre_sign_up"
  role          = aws_iam_role.sonar_service_users_pre_sign_up_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/pre_sign_up"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_post_confirmation" {
  function_name = "sonar_service_users_post_confirmation"
  role          = aws_iam_role.sonar_service_users_post_confirmation_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/post_confirmation"]
  }

  environment {
    variables = {
      "DEFAULT_GROUP" = "externals_guest"
    }
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_external_sign_in" {
  function_name = "sonar_service_users_external_sign_in"
  role          = aws_iam_role.sonar_service_users_external_sign_in_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  timeout       = 5

  environment {
    variables = {
      CLIENT_ID    = aws_cognito_user_pool_client.externals.id
      USER_POOL_ID = aws_cognito_user_pool.externals.id
    }
  }

  image_config {
    entry_point = [
      "/lambda/external_sign_in"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_external_sign_up" {
  function_name = "sonar_service_users_external_sign_up"
  role          = aws_iam_role.sonar_service_users_external_sign_up_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  timeout       = 5

  environment {
    variables = {
      CLIENT_ID = aws_cognito_user_pool_client.externals.id
    }
  }

  image_config {
    entry_point = [
      "/lambda/external_sign_up"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_external_refresh" {
  function_name = "sonar_service_users_external_refresh"
  role          = aws_iam_role.sonar_service_users_external_refresh_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  environment {
    variables = {
      CLIENT_ID = aws_cognito_user_pool_client.externals.id
    }
  }

  image_config {
    entry_point = [
      "/lambda/external_refresh"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}


resource "aws_lambda_function" "sonar_service_users_olive_authorizer" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_users_olive_authorizer"
  role          = aws_iam_role.sonar_service_users_olive_authorizer_role.arn
  image_config {
    entry_point = ["/lambda/authorizer"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_define_auth_challenge" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_users_define_auth_challenge"
  role          = aws_iam_role.sonar_service_users_define_auth_challenge_role.arn
  image_config {
    entry_point = ["/lambda/define_auth_challenge"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_internal_post_authentication" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_users_internal_post_authentication"
  role          = aws_iam_role.sonar_service_users_internal_post_authentication_role.arn
  image_config {
    entry_point = ["/lambda/internal_post_authentication"]
  }

  environment {
    variables = {
      ENVIRONMENT = var.environment
    }
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_internal_post_confirmation" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_users_internal_post_confirmation"
  role          = aws_iam_role.sonar_service_users_internal_post_confirmation_role.arn
  image_config {
    entry_point = ["/lambda/internal_post_confirmation"]
  }

  environment {
    variables = {
      ENVIRONMENT = var.environment
    }
  }

  tracing_config {
    mode = "Active"
  }
}


resource "aws_lambda_function" "sonar_service_users_user_list" {
  function_name = "sonar_service_users_user_list"
  role          = aws_iam_role.sonar_service_users_user_list_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  timeout = 30

  environment {
    variables = {
      USER_POOL_ID = aws_cognito_user_pool.externals.id
      HOST         = var.db_host
      NAME         = var.db_name
      PASSWORD     = var.db_password
      PORT         = var.db_port
      USER         = var.db_username
    }
  }

  image_config {
    entry_point = [
      "/lambda/user_list"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_users_revoke_access" {
  function_name = "sonar_service_users_revoke_access"
  role          = aws_iam_role.sonar_service_users_revoke_access_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  environment {
    variables = {
      USER_POOL_ID  = aws_cognito_user_pool.externals.id
      WEBSOCKET_URL = "https://ws-sonar.${var.domain_name}"
    }
  }

  image_config {
    entry_point = [
      "/lambda/revoke_access"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_users_get_organizations" {
  function_name = "sonar_service_users_get_organizations"
  role          = aws_iam_role.sonar_service_users_get_organizations_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  timeout = 30

  environment {
    variables = {
      USER_POOL_ID = aws_cognito_user_pool.externals.id
      HOST         = var.db_host
      NAME         = var.db_name
      PASSWORD     = var.db_password
      PORT         = var.db_port
      USER         = var.db_username
    }
  }

  image_config {
    entry_point = [
      "/lambda/get_organizations"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_update_user" {
  function_name = "sonar_service_update_user"
  role          = aws_iam_role.sonar_service_update_user_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  timeout       = 10

  environment {
    variables = {
      USER_POOL_ID  = aws_cognito_user_pool.externals.id
      WEBSOCKET_URL = "https://ws-unconfirmed.${var.domain_name}"
      HOST          = var.db_host
      NAME          = var.db_name
      PASSWORD      = var.db_password
      PORT          = var.db_port
      USER          = var.db_username
    }
  }

  image_config {
    entry_point = [
      "/lambda/update_user"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

### task_get_external_user start ###
resource "aws_lambda_function" "sonar_service_users_task_get_external_user" {
  function_name = "sonar_service_users_task_get_external_user"
  role          = aws_iam_role.sonar_service_users_task_get_external_user_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  timeout = 30

  environment {
    variables = {
      USER_POOL_ID = aws_cognito_user_pool.externals.id
      HOST         = var.db_host
      NAME         = var.db_name
      PASSWORD     = var.db_password
      PORT         = var.db_port
      USER         = var.db_username
    }
  }

  image_config {
    entry_point = [
      "/lambda/task_get_external_user"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}
### task_get_external_user end ###

### task_get_internal_user start ###
resource "aws_lambda_function" "sonar_service_users_task_get_internal_user" {
  function_name = "sonar_service_users_task_get_internal_user"
  role          = aws_iam_role.sonar_service_users_task_get_internal_user_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  environment {
    variables = {
      USER_POOL_ID = aws_cognito_user_pool.internals.id
    }
  }

  image_config {
    entry_point = [
      "/lambda/task_get_internal_user"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}
### task_get_internal_user end ###

resource "aws_lambda_function" "sonar_service_users_create_organizations" {
  function_name = "sonar_service_users_create_organizations"
  role          = aws_iam_role.sonar_service_users_create_organizations_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  timeout = 30

  environment {
    variables = {
      USER_POOL_ID = aws_cognito_user_pool.externals.id
      HOST         = var.db_host
      NAME         = var.db_name
      PASSWORD     = var.db_password
      PORT         = var.db_port
      USER         = var.db_username
    }
  }

  image_config {
    entry_point = [
      "/lambda/create_organizations"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}
