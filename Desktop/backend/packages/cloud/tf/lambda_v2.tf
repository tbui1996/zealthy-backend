// Why these are V2: in order to introduce new functionality to loop without breaking it, we sometimes will need to duplicate lambdas.
// These will then use new apigateways to hit new routes on a v2 config so that we can isolate change better.
resource "aws_lambda_function" "sonar_service_cloud_file_upload_loop_v2" {
  function_name = "sonar_service_cloud_file_upload_loop_v2"
  role          = aws_iam_role.sonar_service_cloud_upload_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/file_upload"]
  }

  environment {
    variables = {
      "BUCKETNAME" = aws_s3_bucket.sonar_cloud.id
      "HOST"       = var.db_host
      "NAME"       = var.db_name
      "PASSWORD"   = var.db_password
      "PORT"       = var.db_port
      "USER"       = var.db_username
    }
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

resource "aws_lambda_function" "sonar_service_cloud_file_download_loop_v2" {
  function_name = "sonar_service_cloud_file_download_loop_v2"
  role          = aws_iam_role.sonar_service_cloud_download_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/file_download"]
  }

  environment {
    variables = {
      "BUCKETNAME" = aws_s3_bucket.sonar_cloud.id
    }
  }
  tracing_config {
    mode = "Active"
  }
}

resource "aws_lambda_function" "sonar_service_cloud_pre_signed_upload_url_loop" {
  function_name = "sonar_service_cloud_pre_signed_upload_url_loop_v2"
  role          = aws_iam_role.sonar_service_cloud_pre_signed_upload_url.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/pre_signed_upload_url"]
  }

  environment {
    variables = {
      "BUCKETNAME" = aws_s3_bucket.sonar_cloud.id
    }
  }

  tracing_config {
    mode = "Active"
  }
}