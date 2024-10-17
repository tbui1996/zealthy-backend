resource "aws_iam_role" "aws_apigatewayv2_internal_authorizer_role" {
  name               = "sonar_global_api_gateway_invocation_role"
  assume_role_policy = data.aws_iam_policy_document.apigw_execution_document.json
}

resource "aws_iam_role_policy" "aws_apigatewayv2_internal_authorizer_policy" {
  name = "sonar_global_api_gateway_invocation_policy"
  role = aws_iam_role.aws_apigatewayv2_internal_authorizer_role.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : "lambda:InvokeFunction",
          Effect : "Allow",
          Resource : "${aws_lambda_function.internal_authorizer_lambda.arn}"
        }
      ]
  })
}