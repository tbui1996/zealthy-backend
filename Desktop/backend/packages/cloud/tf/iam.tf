data "aws_iam_policy_document" "log_group_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/sonar_service_cloud_*:*"
    ]
  }
}

data "aws_iam_policy_document" "kms_key_policy_document" {
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
      aws_kms_key.sonar_cloud_s3_bucket_kms_key.arn
    ]
  }
}

data "aws_iam_policy_document" "s3_get_objects_policy_document" {
  statement {
    actions = [
      "s3:GetObject",
      "s3:GetBucketLocation",
      "s3:ListBucket"
    ]

    resources = [
      aws_s3_bucket.sonar_cloud.arn,
      "${aws_s3_bucket.sonar_cloud.arn}/*"
    ]

    condition {
      test     = "StringEquals"
      values   = [data.aws_caller_identity.current.account_id]
      variable = "aws:PrincipalAccount"
    }
  }
}

data "aws_iam_policy_document" "s3_upload_objects_policy_document" {
  statement {
    actions = [
      "s3:AbortMultipartUpload",
      "s3:PutObject",
      "s3:GetBucketLocation",
      "s3:DeleteObject"
    ]

    resources = [
      aws_s3_bucket.sonar_cloud.arn,
      "${aws_s3_bucket.sonar_cloud.arn}/*"
    ]

    condition {
      test     = "StringEquals"
      values   = [data.aws_caller_identity.current.account_id]
      variable = "aws:PrincipalAccount"
    }
  }
}

# -- Associate File Start --
resource "aws_iam_role" "sonar_service_cloud_associate_file" {
  name               = "sonar_service_cloud_associate_file"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_cloud_associate_file_vpc" {
  role       = aws_iam_role.sonar_service_cloud_associate_file.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_associate_file_policy" {
  name   = "sonar_service_cloud_associate_file_policy"
  role   = aws_iam_role.sonar_service_cloud_associate_file.id
  policy = data.aws_iam_policy_document.log_group_policy_document.json
}
# -- Associate File End --

# -- Download File Start --
resource "aws_iam_role" "sonar_service_cloud_download_file" {
  name               = "sonar_service_cloud_download_file"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_cloud_download_file_vpc" {
  role       = aws_iam_role.sonar_service_cloud_download_file.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "kms_log_s3_get_combined_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.log_group_policy_document.json,
    data.aws_iam_policy_document.kms_key_policy_document.json,
    data.aws_iam_policy_document.s3_get_objects_policy_document.json
  ]
}

resource "aws_iam_role_policy" "sonar_service_cloud_download_file_policy" {
  name   = "sonar_service_cloud_download_file_policy"
  role   = aws_iam_role.sonar_service_cloud_download_file.id
  policy = data.aws_iam_policy_document.kms_log_s3_get_combined_policy_document.json
}
# -- Download File End --

# -- Upload File Start --
resource "aws_iam_role" "sonar_service_cloud_upload_file" {
  name               = "sonar_service_cloud_upload_file"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_cloud_upload_file_vpc" {
  role       = aws_iam_role.sonar_service_cloud_upload_file.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "kms_log_s3_upload_combined_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.log_group_policy_document.json,
    data.aws_iam_policy_document.kms_key_policy_document.json,
    data.aws_iam_policy_document.s3_upload_objects_policy_document.json
  ]
}

resource "aws_iam_role_policy" "sonar_service_cloud_upload_file_policy" {
  name   = "sonar_service_cloud_upload_file_policy"
  role   = aws_iam_role.sonar_service_cloud_upload_file.id
  policy = data.aws_iam_policy_document.kms_log_s3_upload_combined_policy_document.json
}
# -- Upload File End --

# -- Get File Start --
resource "aws_iam_role" "sonar_service_cloud_get_file" {
  name               = "sonar_service_cloud_get_file"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_cloud_get_file_vpc" {
  role       = aws_iam_role.sonar_service_cloud_get_file.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_cloud_get_file_policy" {
  name   = "sonar_service_cloud_get_file_policy"
  policy = data.aws_iam_policy_document.log_group_policy_document.json
  role   = aws_iam_role.sonar_service_cloud_get_file.id
}
# -- Get File End --

# -- View File Start --
resource "aws_iam_role" "sonar_service_cloud_view_file" {
  name               = "sonar_service_cloud_view_file"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy" "sonar_service_cloud_view_file_policy" {
  name   = "sonar_service_cloud_view_file_policy"
  policy = data.aws_iam_policy_document.kms_log_s3_get_combined_policy_document.json
  role   = aws_iam_role.sonar_service_cloud_view_file.id
}
## -- View File End --

# -- Delete File Start --
resource "aws_iam_role" "sonar_service_cloud_delete_file" {
  name               = "sonar_service_cloud_delete_file"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_cloud_delete_file_vpc" {
  role       = aws_iam_role.sonar_service_cloud_delete_file.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_cloud_delete_file_policy" {
  name   = "sonar_service_cloud_delete_file_policy"
  role   = aws_iam_role.sonar_service_cloud_delete_file.id
  policy = data.aws_iam_policy_document.kms_log_s3_upload_combined_policy_document.json
}
# -- Delete File End --

# -- Pre-signed upload url Start --
resource "aws_iam_role" "sonar_service_cloud_pre_signed_upload_url" {
  name               = "sonar_service_cloud_pre_signed_upload_url"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

data "aws_iam_policy_document" "kms_log_s3_pre_signed_upload_url_policy_document" {
  source_policy_documents = [
    data.aws_iam_policy_document.log_group_policy_document.json,
    data.aws_iam_policy_document.kms_key_policy_document.json,
    data.aws_iam_policy_document.s3_upload_objects_policy_document.json,
    data.aws_iam_policy_document.s3_get_objects_policy_document.json
  ]
}

resource "aws_iam_role_policy" "sonar_service_cloud_pre_signed_upload_url_policy" {
  name   = "sonar_service_cloud_pre_signed_upload_url_policy"
  role   = aws_iam_role.sonar_service_cloud_pre_signed_upload_url.id
  policy = data.aws_iam_policy_document.kms_log_s3_pre_signed_upload_url_policy_document.json
}
# -- Pre-signed upload url  End --
