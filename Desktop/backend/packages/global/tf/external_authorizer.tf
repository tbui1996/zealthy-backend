/* Authorizer */
resource "aws_iam_role" "external_authorizer_credentials_role" {
  name               = "sonar_service_${var.name}_external_authorizer_credentials_role"
  assume_role_policy = data.aws_iam_policy_document.apigw_execution_document.json
}

resource "aws_iam_role_policy" "external_authorizer_credentials_invoke" {
  name = "external_authorizer_credentials_invoke"
  role = aws_iam_role.external_authorizer_credentials_role.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : [
            "lambda:InvokeFunction"
          ],
          Resource : "${aws_lambda_function.external_authorizer_lambda.arn}",
          Effect : "Allow"
        }
      ]
  })
}

/* Lambda */
resource "aws_iam_role" "external_authorizer_role" {
  name               = "sonar_service_${var.name}_external_authorizer_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "external_authorizer_logging" {
  name = "logging_policy"
  role = aws_iam_role.external_authorizer_role.id
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
          Resource : "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.external_authorizer_lambda.function_name}:*",
          Effect : "Allow"
        }
      ]
  })
}

resource "aws_iam_role_policy" "external_authorizer_dynamodb" {
  name = "dynamodb_policy"
  role = aws_iam_role.external_authorizer_role.id
  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Action : [
            "dynamodb:GetItem",
          ],
          Resource : "arn:aws:dynamodb:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:table/SonarGroupPolicy"
          Effect : "Allow"
        },
        {
          Effect : "Allow"
          Action : [
            "kms:DescribeKey",
            "kms:Encrypt",
            "kms:Decrypt",
            "kms:ReEncryptTo",
            "kms:ReEncryptFrom",
            "kms:GenerateDataKey",
            "kms:GenerateDataKeyWithoutPlaintext"
          ],
          Resource : aws_kms_key.sonar_group_policy_dynamodb_kms_key.arn
        },
      ]
  })
}
