## -- SERVICE_NAME LAMBDA_NAME Start --
resource "aws_iam_role" "sonar_service_SERVICE_NAME_LAMBDA_NAME_role" {
  name               = "sonar_service_SERVICE_NAME_LAMBDA_NAME_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

# resource "aws_iam_role_policy_attachment" "sonar_service_SERVICE_NAME_LAMBDA_NAME_vpc" {
#   role       = aws_iam_role.sonar_service_SERVICE_NAME_LAMBDA_NAME_role.name
#   policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
# }

data "aws_iam_policy_document" "sonar_service_SERVICE_NAME_LAMBDA_NAME_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.apigateway_web_invoke_api_policy_document.json
  ]

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_SERVICE_NAME_LAMBDA_NAME.function_name}*:*",
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_SERVICE_NAME_LAMBDA_NAME_policy" {
  name   = "sonar_service_SERVICE_NAME_LAMBDA_NAME_policy"
  role   = aws_iam_role.sonar_service_SERVICE_NAME_LAMBDA_NAME_role.id
  policy = data.aws_iam_policy_document.sonar_service_SERVICE_NAME_LAMBDA_NAME_policy_document.json
}
## -- SERVICE_NAME LAMBDA_NAME End --