resource "aws_lambda_function" "sonar_service_support_loop_chat_sessions_get" {
  function_name = "sonar_service_support_loop_chat_sessions_get"
  role          = aws_iam_role.sonar_service_support_chat_sessions_get.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/loop_chat_sessions_get"
    ]
  }

  tracing_config {
    mode = "Active"
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
}