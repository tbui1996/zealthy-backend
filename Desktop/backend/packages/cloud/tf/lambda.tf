// Loop lambdas
// TODO: DEPRECATED: This is moving in favor of V2, delete as soon as possible!
resource "aws_lambda_function" "sonar_service_cloud_file_upload_loop" {
  function_name = "sonar_service_cloud_file_upload_loop"
  role          = aws_iam_role.sonar_service_cloud_upload_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.last_production_commit}"

  image_config {
    entry_point = [
    "/lambda/file_upload"]
  }

  environment {
    variables = {
      "BUCKETNAME"   = aws_s3_bucket.sonar_cloud.id
      "MIMEFILEPATH" = "/lambda/mime.json"
      "HOST"         = var.db_host
      "NAME"         = var.db_name
      "PASSWORD"     = var.db_password
      "PORT"         = var.db_port
      "USER"         = var.db_username
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

// TODO: DEPRECATED: This is moving in favor of V2, delete as soon as possible!
resource "aws_lambda_function" "sonar_service_cloud_file_download_loop" {
  function_name = "sonar_service_cloud_file_download_loop"
  role          = aws_iam_role.sonar_service_cloud_download_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.last_production_commit}"

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

// Internal/Web app lambdas
resource "aws_lambda_function" "sonar_service_cloud_file_upload_web" {
  function_name = "sonar_service_cloud_file_upload_web"
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

resource "aws_lambda_function" "sonar_service_cloud_file_download_web" {
  function_name = "sonar_service_cloud_file_download_web"
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

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_cloud_get_file" {
  function_name = "sonar_service_cloud_get_file"
  role          = aws_iam_role.sonar_service_cloud_get_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/get_file"]
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
    subnet_ids         = var.private_subnets
    security_group_ids = [var.rds_security_group]
  }
}

resource "aws_lambda_function" "sonar_service_cloud_associate_file_live_env" {
  count         = var.live_env ? 1 : 0
  function_name = "sonar_service_cloud_associate_file"
  role          = aws_iam_role.sonar_service_cloud_associate_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/associate_file"]
  }

  environment {
    variables = {
      "BUCKETNAME"      = aws_s3_bucket.sonar_cloud.id
      "HOST"            = var.db_host
      "NAME"            = var.db_name
      "PASSWORD"        = var.db_password
      "PORT"            = var.db_port
      "USER"            = var.db_username
      "DOPPLERHOST"     = var.doppler_host
      "DOPPLERPORT"     = var.doppler_port
      "DOPPLERUSER"     = var.doppler_user
      "DOPPLERPASSWORD" = var.doppler_pw
      "DOPPLERDBNAME"   = var.doppler_dbname
    }
  }
  tracing_config {
    mode = "Active"
  }
  vpc_config {
    subnet_ids         = var.private_subnets
    security_group_ids = [var.rds_security_group]
  }
}

resource "aws_lambda_function" "sonar_service_cloud_associate_file_dev_env" {
  count         = var.live_env ? 0 : 1
  function_name = "sonar_service_cloud_associate_file"
  role          = aws_iam_role.sonar_service_cloud_associate_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/associate_file"]
  }

  environment {
    variables = {
      "BUCKETNAME"      = aws_s3_bucket.sonar_cloud.id
      "HOST"            = var.db_host
      "NAME"            = var.db_name
      "PASSWORD"        = var.db_password
      "PORT"            = var.db_port
      "USER"            = var.db_username
      "DOPPLERHOST"     = var.db_host
      "DOPPLERPORT"     = var.db_port
      "DOPPLERUSER"     = var.db_username
      "DOPPLERPASSWORD" = var.db_password
      "DOPPLERDBNAME"   = "external"
    }
  }
  tracing_config {
    mode = "Active"
  }
  vpc_config {
    subnet_ids         = var.private_subnets
    security_group_ids = [var.rds_security_group]
  }
}

resource "aws_lambda_function" "sonar_service_cloud_delete_file" {
  function_name = "sonar_service_cloud_delete_file"
  role          = aws_iam_role.sonar_service_cloud_delete_file.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/delete_file"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_cloud_pre_signed_upload_url_web" {
  function_name = "sonar_service_cloud_pre_signed_upload_url_web"
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
