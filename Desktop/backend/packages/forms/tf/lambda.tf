resource "aws_lambda_function" "sonar_service_forms_get" {
  function_name = "sonar_service_forms_get"
  role          = aws_iam_role.sonar_service_form_get.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/get"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_list" {
  function_name = "sonar_service_forms_list"
  role          = aws_iam_role.sonar_service_form_list.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/list"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_create" {
  function_name = "sonar_service_forms_create"
  role          = aws_iam_role.sonar_service_form_create.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/create"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_send" {
  function_name = "sonar_service_forms_send"
  role          = aws_iam_role.sonar_service_form_send.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/send"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
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

resource "aws_lambda_function" "sonar_service_forms_response" {
  function_name = "sonar_service_forms_response"
  role          = aws_iam_role.sonar_service_form_response.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/response"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_count" {
  function_name = "sonar_service_forms_count"
  role          = aws_iam_role.sonar_service_form_count.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/count"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_receive" {
  function_name = "sonar_service_forms_receive"
  role          = aws_iam_role.sonar_service_form_receive.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/receive"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
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

resource "aws_lambda_function" "sonar_service_forms_delete" {
  function_name = "sonar_service_forms_delete"
  role          = aws_iam_role.sonar_service_form_delete.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/delete"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_edit" {
  function_name = "sonar_service_forms_edit"
  role          = aws_iam_role.sonar_service_form_edit.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/edit"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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

resource "aws_lambda_function" "sonar_service_forms_close" {
  function_name = "sonar_service_forms_close"
  role          = aws_iam_role.sonar_service_form_close.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/close"]
  }

  environment {
    variables = {
      "HOST"     = var.db_host
      "NAME"     = var.form_db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
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
