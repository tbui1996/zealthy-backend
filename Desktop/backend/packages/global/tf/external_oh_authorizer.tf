resource "aws_iam_role" "external_oh_authorizer" {
  name               = "external_oh_authorizer"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "external_oh_authorizer" {
  name = "external_oh_authorizer"
  role = aws_iam_role.external_oh_authorizer.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ],
          Resource : "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.external_oh_authorizer_lambda.function_name}:*",
          Effect : "Allow"
        }
      ]
  })
}


resource "aws_iam_role" "external_oh_authorizer_credentials_role" {
  name               = "global_external_oh_authorizer_credentials_role"
  assume_role_policy = data.aws_iam_policy_document.apigw_execution_document.json
}

resource "aws_iam_role_policy" "external_oh_authorizer_credentials_invoke" {
  name = "global_external_oh_authorizer_credentials_invoke"
  role = aws_iam_role.external_oh_authorizer_credentials_role.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : [
            "lambda:InvokeFunction"
          ],
          Resource : "${aws_lambda_function.external_oh_authorizer_lambda.arn}",
          Effect : "Allow"
        }
      ]
  })
}
