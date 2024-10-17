### All one shot lambda config for syncing cognito users are here so that it can be moved/removed later

###################### Lambda ########################
resource "aws_lambda_function" "sonar_service_users_sync_cognito_users" {
  function_name = "sonar_service_users_sync_cognito_users"
  role          = aws_iam_role.sonar_service_users_sync_cognito_users_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  timeout       = 30

  environment {
    variables = {
      HOST         = var.db_host
      NAME         = var.db_name
      PASSWORD     = var.db_password
      PORT         = var.db_port
      USER         = var.db_username
      USER_POOL_ID = aws_cognito_user_pool.externals.id
    }
  }

  image_config {
    entry_point = [
      "/lambda/sync_cognito_users"
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

###################### IAM ########################
#### sync_cognito_users lambda start
resource "aws_iam_role" "sonar_service_users_sync_cognito_users_role" {
  name               = "sonar_service_users_sync_cognito_users_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_users_sync_cognito_users_policy" {
  name   = "sonar_service_users_sync_cognito_users_policy"
  role   = aws_iam_role.sonar_service_users_sync_cognito_users_role.id
  policy = data.aws_iam_policy_document.sonar_service_users_sync_cognito_users_policy_document.json
}

data "aws_iam_policy_document" "sonar_service_users_sync_cognito_users_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_users_sync_cognito_users.function_name}:*"]
  }

  statement {
    actions = [
      "cognito-idp:ListUsers",
    ]
    resources = [aws_cognito_user_pool.externals.arn]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_users_sync_cognito_lambda_vpc" {
  role       = aws_iam_role.sonar_service_users_sync_cognito_users_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
#### sync_cognito_users lambda end
