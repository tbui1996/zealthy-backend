resource "aws_lambda_function" "internal_authorizer_lambda" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_global_internal_authorizer"
  role          = aws_iam_role.internal_authorizer_role.arn
  image_config {
    entry_point = ["/lambda/sonar_authorizer"]
  }

  environment {
    variables = {
      JWKS_URL          = var.aws_cognito_user_pool_internals_jwks
      GROUP_NAME_PREFIX = "internals"
    }
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "external_authorizer_lambda" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_global_external_authorizer"
  role          = aws_iam_role.external_authorizer_role.arn
  image_config {
    entry_point = ["/lambda/sonar_authorizer"]
  }

  environment {
    variables = {
      JWKS_URL          = var.aws_cognito_user_pool_externals_jwks
      GROUP_NAME_PREFIX = "externals"
    }
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "external_oh_authorizer_lambda" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_global_external_oh_authorizer"
  role          = aws_iam_role.external_oh_authorizer.arn

  image_config {
    entry_point = ["/lambda/external_oh_authorizer"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "internal_websocket_connect" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_global_internal_websocket_connect"
  role          = aws_iam_role.global_route_connect_role.arn

  image_config {
    entry_point = ["/lambda/connect"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "internal_websocket_disconnect" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_global_internal_websocket_disconnect"
  role          = aws_iam_role.global_route_disconnect_role.arn

  image_config {
    entry_point = ["/lambda/disconnect"]
  }

  tracing_config {
    mode = "Active"
  }
}

### Start Is Internal User Online ###
resource "aws_lambda_function" "global_task_is_internal_user_online" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "global_task_is_internal_user_online"
  role          = aws_iam_role.global_task_is_internal_user_online_role.arn
  image_config {
    entry_point = ["/lambda/task_is_internal_user_online"]
  }

  tracing_config {
    mode = "Active"
  }
}
### End Is Internal User Online