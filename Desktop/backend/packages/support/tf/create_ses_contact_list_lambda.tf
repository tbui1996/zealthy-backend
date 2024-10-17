### All one shot lambda config for creating a AWS SES Contact list here so that it can be moved/removed later

###################### Lambda ########################
resource "aws_lambda_function" "create_ses_contact_list" {
  function_name = "create_ses_contact_list"
  role          = aws_iam_role.create_ses_contact_list_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/create_ses_contact_list"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}

###################### IAM ########################
### Create SES Contact List Permissions ###
resource "aws_iam_role" "create_ses_contact_list_role" {
  name               = "create_ses_contact_list_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "create_ses_contact_list_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
  ]

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
      aws_kms_key.offline_messages_kms_key.arn,
    ]
  }

  statement {
    actions = [
      "ses:CreateContactList",
    ]
    resources = ["arn:aws:ses:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"]
  }
}

resource "aws_iam_role_policy" "create_ses_contact_list_policy" {
  name = "create_ses_contact_list_policy"
  role = aws_iam_role.create_ses_contact_list_role.id

  policy = data.aws_iam_policy_document.create_ses_contact_list_document.json
}
### End Create SES Contact List Permissions ###