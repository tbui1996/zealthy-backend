resource "aws_lambda_function" "sonar_service_feature_flags_create_flag" {
  function_name = "sonar_service_feature_flags_create_flag"
  role          = aws_iam_role.sonar_service_feature_flags_create_flag_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/create_flag"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
      "HOST"     = var.db_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_feature_flags_patch_flag" {
  function_name = "sonar_service_feature_flags_patch_flag"
  role          = aws_iam_role.sonar_service_feature_flags_patch_flag_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/patch_flag"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
      "HOST"     = var.db_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}


resource "aws_lambda_function" "sonar_service_feature_flags_list_flags" {
  function_name = "sonar_service_feature_flags_list_flags"
  role          = aws_iam_role.sonar_service_feature_flags_list_flags_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/list_flags"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
      "HOST"     = var.db_read_only_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_feature_flags_evaluate" {
  function_name = "sonar_service_feature_flags_evaluate"
  role          = aws_iam_role.sonar_service_feature_flags_evaluate_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/evaluate"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
      "HOST"     = var.db_read_only_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_feature_flags_loop_evaluate" {
  function_name = "sonar_service_feature_flags_loop_evaluate"
  role          = aws_iam_role.sonar_service_feature_flags_loop_evaluate_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/evaluate"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
      "HOST"     = var.db_read_only_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}

resource "aws_lambda_function" "sonar_service_feature_flags_delete_flag" {
  function_name = "sonar_service_feature_flags_delete_flag"
  role          = aws_iam_role.sonar_service_feature_flags_delete_flag_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/delete_flag"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
      "HOST"     = var.db_host
      "NAME"     = var.db_name
      "PASSWORD" = var.db_password
      "PORT"     = var.db_port
      "USER"     = var.db_username
    }
  }

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}