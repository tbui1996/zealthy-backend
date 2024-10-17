## -- patient list_appointments begin -- 

resource "aws_iam_role" "sonar_service_patient_list_appointments_role" {
  name               = "sonar_service_patient_list_appointments_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_list_appointments_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_list_appointments.function_name}:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/GET/appointments"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_list_appointments_vpc" {
  role       = aws_iam_role.sonar_service_patient_list_appointments_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_list_appointments_policy" {
  name   = "sonar_service_patient_list_appointments_policy"
  role   = aws_iam_role.sonar_service_patient_list_appointments_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_list_appointments_policy_document.json
}
## -- patient list_appointments End -- 

## -- patient create_appointments begin --
resource "aws_iam_role" "sonar_service_patient_create_appointments_role" {
  name               = "sonar_service_patient_create_appointments_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_create_appointments_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_create_appointments.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/POST/appointments"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_create_appointments_vpc" {
  role       = aws_iam_role.sonar_service_patient_create_appointments_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_create_appointments_policy" {
  name   = "sonar_service_patient_create_appointments_policy"
  role   = aws_iam_role.sonar_service_patient_create_appointments_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_create_appointments_policy_document.json
}

## -- patient create_appointments End

## -- patient edit_appointments begin --
resource "aws_iam_role" "sonar_service_patient_edit_appointments_role" {
  name               = "sonar_service_patient_edit_appointments_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_edit_appointments_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_edit_appointments.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/PUT/appointments/*"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_edit_appointments_vpc" {
  role       = aws_iam_role.sonar_service_patient_edit_appointments_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_edit_appointments_policy" {
  name   = "sonar_service_patient_edit_appointments_policy"
  role   = aws_iam_role.sonar_service_patient_edit_appointments_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_edit_appointments_policy_document.json
}

## -- patient edit_appointments End

## -- patient delete_appointments begin --
resource "aws_iam_role" "sonar_service_patient_delete_appointments_role" {
  name               = "sonar_service_patient_delete_appointments_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_delete_appointments_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_delete_appointments.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/DELETE/appointments/*"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_delete_appointments_vpc" {
  role       = aws_iam_role.sonar_service_patient_delete_appointments_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_delete_appointments_policy" {
  name   = "sonar_service_patient_delete_appointments_policy"
  role   = aws_iam_role.sonar_service_patient_delete_appointments_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_delete_appointments_policy_document.json
}

## -- patient delete_appointments End


## -- patient list_patients begin --
resource "aws_iam_role" "sonar_service_patient_list_patients_role" {
  name               = "sonar_service_patient_list_patients_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_list_patients_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_list_patients.function_name}:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/GET/patients"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_list_patients_vpc" {
  role       = aws_iam_role.sonar_service_patient_list_patients_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_list_patients_policy" {
  name   = "sonar_service_patient_list_patients_policy"
  role   = aws_iam_role.sonar_service_patient_list_patients_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_list_patients_policy_document.json
}

## -- patient list_patients End


## -- patient create_patients begin --
resource "aws_iam_role" "sonar_service_patient_create_patients_role" {
  name               = "sonar_service_patient_create_patients_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_create_patients_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_create_patients.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/POST/patients"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_create_patients_vpc" {
  role       = aws_iam_role.sonar_service_patient_create_patients_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_create_patients_policy" {
  name   = "sonar_service_patient_create_patients_policy"
  role   = aws_iam_role.sonar_service_patient_create_patients_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_create_patients_policy_document.json
}

## -- patient create_patients End

## -- patient patch_patients begin --
resource "aws_iam_role" "sonar_service_patient_patch_patients_role" {
  name               = "sonar_service_patient_patch_patients_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_patch_patients_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_patch_patients.function_name}:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/PUT/patients"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_patch_patients_vpc" {
  role       = aws_iam_role.sonar_service_patient_patch_patients_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_patch_patients_policy" {
  name   = "sonar_service_patient_patch_patients_policy"
  role   = aws_iam_role.sonar_service_patient_patch_patients_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_patch_patients_policy_document.json
}

## -- patient patch_patients End

## -- patient list_agency_providers begin -- 

resource "aws_iam_role" "sonar_service_patient_list_agency_providers_role" {
  name               = "sonar_service_patient_list_agency_providers_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_list_agency_providers_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_list_agency_providers.function_name}:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/GET/agency_providers"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_list_agency_providers_vpc" {
  role       = aws_iam_role.sonar_service_patient_list_agency_providers_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_list_agency_providers_policy" {
  name   = "sonar_service_patient_list_agency_providers_policy"
  role   = aws_iam_role.sonar_service_patient_list_agency_providers_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_list_agency_providers_policy_document.json
}
## -- patient list_agency_providers End -- 

## -- patient create_agency_providers begin --
resource "aws_iam_role" "sonar_service_patient_create_agency_providers_role" {
  name               = "sonar_service_patient_create_agency_providers_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_create_agency_providers_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_create_agency_providers.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/POST/agency_providers"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_create_agency_providers_vpc" {
  role       = aws_iam_role.sonar_service_patient_create_agency_providers_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_create_agency_providers_policy" {
  name   = "sonar_service_patient_create_agency_providers_policy"
  role   = aws_iam_role.sonar_service_patient_create_agency_providers_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_create_agency_providers_policy_document.json
}

## -- patient create_agency_providers End

## -- patient edit_agency_providers begin --
resource "aws_iam_role" "sonar_service_patient_edit_agency_providers_role" {
  name               = "sonar_service_patient_edit_agency_providers_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
}


data "aws_iam_policy_document" "sonar_service_patient_edit_agency_providers_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${aws_lambda_function.sonar_service_patient_edit_agency_providers.function_name}*:*",
    ]
  }

  statement {
    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.apigateway_id}/${var.apigateway_stage_name}/PUT/agency_providers/*"
    ]
  }
}

resource "aws_iam_role_policy_attachment" "sonar_service_patient_edit_agency_providers_vpc" {
  role       = aws_iam_role.sonar_service_patient_edit_agency_providers_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy" "sonar_service_patient_edit_agency_providers_policy" {
  name   = "sonar_service_patient_edit_agency_providers_policy"
  role   = aws_iam_role.sonar_service_patient_edit_agency_providers_role.id
  policy = data.aws_iam_policy_document.sonar_service_patient_edit_agency_providers_policy_document.json
}

## -- patient edit_agency_providers End