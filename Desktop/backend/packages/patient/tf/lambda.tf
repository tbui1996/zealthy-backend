resource "aws_lambda_function" "sonar_service_patient_list_appointments" {
  function_name = "sonar_service_patient_list_appointments"
  role          = aws_iam_role.sonar_service_patient_list_appointments_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/list_appointments"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

resource "aws_lambda_function" "sonar_service_patient_create_appointments" {
  function_name = "sonar_service_patient_create_appointments"
  role          = aws_iam_role.sonar_service_patient_create_appointments_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/create_appointments"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

## -- patient edit_appointments begin --
resource "aws_lambda_function" "sonar_service_patient_edit_appointments" {
  function_name = "sonar_service_patient_edit_appointments"
  role          = aws_iam_role.sonar_service_patient_edit_appointments_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/edit_appointments"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

## -- patient edit_appointments End

## -- patient delete_appointments begin --
resource "aws_lambda_function" "sonar_service_patient_delete_appointments" {
  function_name = "sonar_service_patient_delete_appointments"
  role          = aws_iam_role.sonar_service_patient_delete_appointments_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/delete_appointments"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

## -- patient delete_appointments End

resource "aws_lambda_function" "sonar_service_patient_list_patients" {
  function_name = "sonar_service_patient_list_patients"
  role          = aws_iam_role.sonar_service_patient_list_patients_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/list_patients"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"

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


resource "aws_lambda_function" "sonar_service_patient_create_patients" {
  function_name = "sonar_service_patient_create_patients"
  role          = aws_iam_role.sonar_service_patient_create_patients_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/create_patients"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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
resource "aws_lambda_function" "sonar_service_patient_patch_patients" {
  function_name = "sonar_service_patient_patch_patients"
  role          = aws_iam_role.sonar_service_patient_patch_patients_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/patch_patient"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"

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


resource "aws_lambda_function" "sonar_service_patient_list_agency_providers" {
  function_name = "sonar_service_patient_list_agency_providers"
  role          = aws_iam_role.sonar_service_patient_list_agency_providers_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/list_agency_providers"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

## -- patient create_agency_providers begin --
resource "aws_lambda_function" "sonar_service_patient_create_agency_providers" {
  function_name = "sonar_service_patient_create_agency_providers"
  role          = aws_iam_role.sonar_service_patient_create_agency_providers_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/create_agency_providers"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

## -- patient create_agency_providers End

## -- patient edit_agency_providers begin --
resource "aws_lambda_function" "sonar_service_patient_edit_agency_providers" {
  function_name = "sonar_service_patient_edit_agency_providers"
  role          = aws_iam_role.sonar_service_patient_edit_agency_providers_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/edit_agency_providers"]
  }

  environment {
    variables = {
      "DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
      "DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
      "DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
      "DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
      "DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
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

## -- patient edit_agency_providers End