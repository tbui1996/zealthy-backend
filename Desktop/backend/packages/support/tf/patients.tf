# start patients get loop
resource "aws_lambda_function" "sonar_service_support_loop_patients_get" {
  function_name = "sonar_service_support_loop_patients_get"
  role          = aws_iam_role.sonar_service_support_patients_get.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/patients_get"
    ]
  }

  tracing_config {
    mode = "Active"
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

  vpc_config {
    subnet_ids = var.private_subnets
    security_group_ids = [
      var.rds_security_group,
      var.external_security_group
    ]
  }
}
# end patients get loop

## -- Chat Messages Get Start --
resource "aws_iam_role" "sonar_service_support_patients_get" {
  name               = "sonar_service_support_patients_get"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}

resource "aws_iam_role_policy_attachment" "sonar_service_support_patients_get_vpc" {
  role       = aws_iam_role.sonar_service_support_patients_get.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

data "aws_iam_policy_document" "sonar_service_support_patients_get_policy" {
  source_policy_documents = [
    data.aws_iam_policy_document.logging_policy_document.json,
    data.aws_iam_policy_document.apigateway_loop_invoke_api_policy_document.json
  ]
}

resource "aws_iam_role_policy" "support_patients_get_role_policy" {
  name   = "support_patients_get_policy"
  role   = aws_iam_role.sonar_service_support_patients_get.id
  policy = data.aws_iam_policy_document.sonar_service_support_patients_get_policy.json
}
## -- Chat Messages Get End --