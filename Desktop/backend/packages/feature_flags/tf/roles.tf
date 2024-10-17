## -- feature_flags create_flag Start --
resource "aws_iam_role" "sonar_service_feature_flags_create_flag_role" {
  name               = "sonar_service_feature_flags_create_flag_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_feature_flags_create_flag_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_feature_flags_create_flag.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/POST/feature_flags"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_feature_flags_create_flag_vpc" {
  role       = aws_iam_role.sonar_service_feature_flags_create_flag_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_feature_flags_create_flag_policy" {
  name   = "sonar_service_feature_flags_create_flag_policy"
  role   = aws_iam_role.sonar_service_feature_flags_create_flag_role.id
  policy = data.aws_iam_policy_document.sonar_service_feature_flags_create_flag_policy_document.json
}
## -- feature_flags create_flag End --

## -- feature_flags list_flags Start --
resource "aws_iam_role" "sonar_service_feature_flags_list_flags_role" {
  name               = "sonar_service_feature_flags_list_flags_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_feature_flags_list_flags_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_feature_flags_list_flags.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/GET/feature_flags"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_feature_flags_list_flags_vpc" {
  role       = aws_iam_role.sonar_service_feature_flags_list_flags_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_feature_flags_list_flags_policy" {
  name   = "sonar_service_feature_flags_list_flags_policy"
  role   = aws_iam_role.sonar_service_feature_flags_list_flags_role.id
  policy = data.aws_iam_policy_document.sonar_service_feature_flags_list_flags_policy_document.json
}
## -- feature_flags list_flags End --


## -- feature_flags patch_flag Start --
resource "aws_iam_role" "sonar_service_feature_flags_patch_flag_role" {
  name               = "sonar_service_feature_flags_patch_flag_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_feature_flags_patch_flag_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_feature_flags_patch_flag.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/PATCH/feature_flags"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_feature_flags_patch_flag_vpc" {
  role       = aws_iam_role.sonar_service_feature_flags_patch_flag_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_feature_flags_patch_flag_policy" {
  name   = "sonar_service_feature_flags_patch_flag_policy"
  role   = aws_iam_role.sonar_service_feature_flags_patch_flag_role.id
  policy = data.aws_iam_policy_document.sonar_service_feature_flags_patch_flag_policy_document.json
}
## -- feature_flags patch_flag End --

## -- feature_flags evaluate Start --
resource "aws_iam_role" "sonar_service_feature_flags_evaluate_role" {
  name               = "sonar_service_feature_flags_evaluate_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_feature_flags_evaluate_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_feature_flags_evaluate.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/GET/feature_flags"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_feature_flags_evaluate_vpc" {
  role       = aws_iam_role.sonar_service_feature_flags_evaluate_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_feature_flags_evaluate_policy" {
  name   = "sonar_service_feature_flags_evaluate_policy"
  role   = aws_iam_role.sonar_service_feature_flags_evaluate_role.id
  policy = data.aws_iam_policy_document.sonar_service_feature_flags_evaluate_policy_document.json
}
## -- feature_flags evaluate End --

## -- feature_flags loop_evaluate Start --
resource "aws_iam_role" "sonar_service_feature_flags_loop_evaluate_role" {
  name               = "sonar_service_feature_flags_loop_evaluate_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_feature_flags_loop_evaluate_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_feature_flags_loop_evaluate.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.aws_apigatewayv2_loop_gateway_id}/${var.aws_apigatewayv2_loop_gateway_stage_name}/GET/feature_flags"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_feature_flags_loop_evaluate_vpc" {
  role       = aws_iam_role.sonar_service_feature_flags_loop_evaluate_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_feature_flags_loop_evaluate_policy" {
  name   = "sonar_service_feature_flags_loop_evaluate_policy"
  role   = aws_iam_role.sonar_service_feature_flags_loop_evaluate_role.id
  policy = data.aws_iam_policy_document.sonar_service_feature_flags_loop_evaluate_policy_document.json
}
## -- feature_flags loop_evaluate End --

resource "aws_iam_role" "sonar_service_feature_flags_delete_flag_role" {
  name               = "sonar_service_feature_flags_delete_flag_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "sonar_service_feature_flags_delete_flag_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_feature_flags_delete_flag.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/DELETE/feature_flags"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_feature_flags_delete_flag_vpc" {
  role       = aws_iam_role.sonar_service_feature_flags_delete_flag_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_feature_flags_delete_flag_policy" {
  name   = "sonar_service_feature_flags_delete_flag_policy"
  role   = aws_iam_role.sonar_service_feature_flags_delete_flag_role.id
  policy = data.aws_iam_policy_document.sonar_service_feature_flags_delete_flag_policy_document.json
}
## -- feature_flags delete_flag End --
