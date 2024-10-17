resource "aws_lambda_function" "api_users_list" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_router_users_list"
  role          = aws_iam_role.api_users_list_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/users_list"]
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

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids         = var.private_subnets
    security_group_ids = [var.rds_security_group]
  }
}

resource "aws_lambda_function" "api_broadcast_lambda" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_router_broadcast"
  role          = aws_iam_role.api_broadcast_lambda_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/broadcast"]
  }
  environment {
    variables = {
      API_REGION    = "${data.aws_region.current.name}"
      WEBSOCKET_URL = "https://ws-sonar.${var.domain_name}"
    }
  }

  tracing_config {
    mode = "Active"
  }
}



resource "aws_lambda_function" "sonar_service_router_receive" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_receive"
  role          = aws_iam_role.sonar_service_router_receive_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/receive"]
  }
  environment {
    variables = {
      API_REGION    = "${data.aws_region.current.name}"
      WEBSOCKET_URL = "https://ws-sonar.${var.domain_name}"
    }
  }

  tracing_config {
    mode = "Active"
  }
}

/*** Websocket Route Lambdas ***/
resource "aws_lambda_function" "sonar_service_router_connect" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_connect"
  role          = aws_iam_role.router_route_connect_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/connect"]
  }


  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_router_disconnect" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_disconnect"
  role          = aws_iam_role.router_route_disconnect_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/disconnect"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_router_message" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_message"
  role          = aws_iam_role.router_route_message_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/message"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_router_unconfirmed_connect" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_unconfirmed_connect"
  role          = aws_iam_role.router_unconfirmed_route_connect_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/unconfirmed_connect"]
  }

  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_unconfirmed_disconnect" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_unconfirmed_disconnect"
  role          = aws_iam_role.router_unconfirmed_route_disconnect_role.arn
  timeout       = "29"
  image_config {
    entry_point = ["/lambda/unconfirmed_disconnect"]
  }

  tracing_config {
    mode = "Active"
  }
}

### Start Is External User Online ###
resource "aws_lambda_function" "sonar_service_router_task_is_external_user_online" {
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  function_name = "sonar_service_router_task_is_external_user_online"
  role          = aws_iam_role.router_task_is_external_user_online_role.arn
  image_config {
    entry_point = ["/lambda/task_is_external_user_online"]
  }

  tracing_config {
    mode = "Active"
  }
}
### End Is External User Online
