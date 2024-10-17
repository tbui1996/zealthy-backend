data "aws_iam_policy_document" "logging_role_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/sonar_service_forms_*:*",
    ]
  }
}

## -- Form Count Start --
resource "aws_iam_role" "sonar_service_form_count" {
  name               = "sonar_service_form_count"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_count_vpc" {
  role       = aws_iam_role.sonar_service_form_count.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_count_policy" {
  name   = "sonar_service_form_count_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_count.id
}
## -- Form Count End --

## -- Form Create Start --
resource "aws_iam_role" "sonar_service_form_create" {
  name               = "sonar_service_form_create"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_create_vpc" {
  role       = aws_iam_role.sonar_service_form_create.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_create_policy" {
  name   = "sonar_service_form_create_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_create.id
}
## -- Form Create End --

## -- Form Get Start --
resource "aws_iam_role" "sonar_service_form_get" {
  name               = "sonar_service_form_get"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_get_vpc" {
  role       = aws_iam_role.sonar_service_form_get.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_get_policy" {
  name   = "sonar_service_form_get_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_get.id
}
## -- Form Get End --

## -- Form List Start --
resource "aws_iam_role" "sonar_service_form_list" {
  name               = "sonar_service_form_list"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_list_vpc" {
  role       = aws_iam_role.sonar_service_form_list.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_list_policy" {
  name   = "sonar_service_form_get_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_list.id
}
## -- Form List End --

## -- Form Receive Start --
resource "aws_iam_role" "sonar_service_form_receive" {
  name               = "sonar_service_form_receive"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_receive_vpc" {
  role       = aws_iam_role.sonar_service_form_receive.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "receive_sqs_policy_document" {
  source_json = data.aws_iam_policy_document.logging_role_policy_document.json
  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes"
    ]
    resources = [
      var.receive_queue_arn
    ]
  }
  statement {
    actions = [
      "kms:DescribeKey",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncryptTo",
      "kms:ReEncryptFrom",
      "kms:GenerateDataKey",
      "kms:GenerateDataKeyWithoutPlaintext"
    ]
    resources = [
      var.receive_queue_kms_arn
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_form_receive_policy" {
  name   = "sonar_service_form_receive_policy"
  policy = data.aws_iam_policy_document.receive_sqs_policy_document.json
  role   = aws_iam_role.sonar_service_form_receive.id
}
## -- Form Receive End --

## -- Form Response Start --
resource "aws_iam_role" "sonar_service_form_response" {
  name               = "sonar_service_form_response"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_response_vpc" {
  role       = aws_iam_role.sonar_service_form_response.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_response_policy" {
  name   = "sonar_service_form_response_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_response.id
}
## -- Form Response End --

## -- Form Send Start --
resource "aws_iam_role" "sonar_service_form_send" {
  name               = "sonar_service_form_send"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_send_vpc" {
  role       = aws_iam_role.sonar_service_form_send.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "send_sqs_policy_document" {
  source_json = data.aws_iam_policy_document.logging_role_policy_document.json
  statement {
    actions = [
      "sqs:GetQueueUrl",
      "sqs:SendMessage"
    ]
    resources = [
      var.send_queue_arn
    ]
  }
  statement {
    actions = [
      "kms:DescribeKey",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncryptTo",
      "kms:ReEncryptFrom",
      "kms:GenerateDataKey",
      "kms:GenerateDataKeyWithoutPlaintext"
    ]
    resources = [
      var.send_queue_kms_arn
    ]
  }
}

resource "aws_iam_role_policy" "sonar_service_form_send_policy" {
  name   = "sonar_service_form_send_policy"
  policy = data.aws_iam_policy_document.send_sqs_policy_document.json
  role   = aws_iam_role.sonar_service_form_send.id
}
## -- Form Send End --

## -- Form Delete Start --
resource "aws_iam_role" "sonar_service_form_delete" {
  name               = "sonar_service_form_delete"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_delete_vpc" {
  role       = aws_iam_role.sonar_service_form_delete.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_delete_policy" {
  name   = "sonar_service_form_delete_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_delete.id
}
## -- Form Delete End --

## -- Form Edit Start --
resource "aws_iam_role" "sonar_service_form_edit" {
  name               = "sonar_service_form_edit"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_edit_vpc" {
  role       = aws_iam_role.sonar_service_form_edit.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}


resource "aws_iam_role_policy" "sonar_service_form_edit_policy" {
  name   = "sonar_service_form_edit_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_edit.id
}
## -- Form Edit End --

## -- Form Close Start --
resource "aws_iam_role" "sonar_service_form_close" {
  name               = "sonar_service_form_close"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_form_close_vpc" {
  role       = aws_iam_role.sonar_service_form_close.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_form_close_policy" {
  name   = "sonar_service_form_close_policy"
  policy = data.aws_iam_policy_document.logging_role_policy_document.json
  role   = aws_iam_role.sonar_service_form_close.id
}
## -- Form Close End --
