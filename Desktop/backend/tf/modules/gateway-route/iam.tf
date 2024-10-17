resource "aws_iam_role" "api_gateway_invocation_role" {
  # gir = gateway_invocation_role
  name               = "${var.lambda_function.function_name}-agir"
  assume_role_policy = data.aws_iam_policy_document.apigw_execution_document.json
}

resource "aws_iam_role_policy" "api_gateway_invocation_policy" {
  # gip = gateway_invocation_policy
  name   = "${var.lambda_function.function_name}_gateway_route"
  role   = aws_iam_role.api_gateway_invocation_role.id
  policy = data.aws_iam_policy_document.api_gateway_invocation_policy_document.json
}

data "aws_iam_policy_document" "api_gateway_invocation_policy_document" {
  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = [var.lambda_function.arn]
  }
}
