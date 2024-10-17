resource "aws_lambda_function" "sonar_service_support_receive" {
  function_name = "sonar_service_support_receive"
  role          = aws_iam_role.sonar_service_support_receive.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/receive"]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION    = data.aws_region.current.name
      WEBSOCKET_URL = "https://ws-sonar-internal.${var.domain_name}"
      "HOST"        = var.db_host
      "NAME"        = var.db_name
      "PASSWORD"    = var.db_password
      "PORT"        = var.db_port
      "USER"        = var.db_username
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

resource "aws_lambda_function" "sonar_service_support_chat_session_create" {
  function_name = "sonar_service_support_chat_session_create"
  role          = aws_iam_role.support_chat_session_create.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/chat_session_create"]
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

resource "aws_lambda_function" "sonar_service_support_pending_chat_session_create" {
  function_name = "sonar_service_support_pending_chat_session_create"
  role          = aws_iam_role.support_pending_chat_session_create.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/pending_chat_session_create"]
  }

  environment {
    variables = {
      API_REGION    = data.aws_region.current.name
      WEBSOCKET_URL = "https://ws-sonar-internal.${var.domain_name}"
      "HOST"        = var.db_host
      "NAME"        = var.db_name
      "PASSWORD"    = var.db_password
      "PORT"        = var.db_port
      "USER"        = var.db_username
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

resource "aws_lambda_function" "sonar_service_support_pending_chat_sessions_get" {
  function_name = "sonar_service_support_pending_chat_sessions_get"
  role          = aws_iam_role.support_pending_chat_sessions_get.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/pending_chat_sessions_get"]
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

resource "aws_lambda_function" "sonar_service_support_assign_pending_chat_session" {
  function_name = "sonar_service_support_assign_pending_chat_session"
  role          = aws_iam_role.sonar_service_support_assign_pending_chat_session.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
    "/lambda/assign_pending_chat_session"]
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

resource "aws_lambda_function" "sonar_service_support_chat_messages_get" {
  function_name = "sonar_service_support_chat_messages_get"
  role          = aws_iam_role.sonar_service_support_chat_messages_get.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/chat_messages_get"
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

resource "aws_lambda_function" "sonar_service_support_chat_sessions_get" {
  function_name = "sonar_service_support_chat_sessions_get"
  role          = aws_iam_role.sonar_service_support_chat_sessions_get.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/chat_sessions_get"
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

resource "aws_lambda_function" "sonar_service_support_chat_session_update_open" {
  function_name = "sonar_service_support_chat_session_update_status"
  role          = aws_iam_role.sonar_service_support_chat_session_update.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/chat_session_update_status"
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

resource "aws_lambda_function" "sonar_service_support_send" {
  function_name = "sonar_service_support_send"
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  role          = aws_iam_role.sonar_service_support_send.arn
  timeout       = 29

  image_config {
    entry_point = ["/lambda/send"]
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

resource "aws_lambda_function" "sonar_service_support_update_chat_notes" {
  function_name = "sonar_service_support_update_chat_notes"
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  role          = aws_iam_role.sonar_service_support_update_chat_notes.arn
  timeout       = 29

  image_config {
    entry_point = ["/lambda/update_chat_notes"]
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

resource "aws_lambda_function" "sonar_service_support_submit_feedback" {
  function_name = "sonar_service_support_submit_feedback"
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  role          = aws_iam_role.sonar_service_support_submit_feedback.arn
  timeout       = 29

  image_config {
    entry_point = ["/lambda/submit_feedback"]
  }

  environment {
    variables = {
      "EMAIL_TEMPLATE"         = var.email_template_feedback
      "CONFIGURATION_SET_NAME" = var.configuration_set_name
      "EMAIL_IDENTITY"         = var.email_identity
    }
  }

  tracing_config {
    mode = "Active"
  }
}

### Start consume event state handler ###
resource "aws_lambda_function" "offline_message_notifier_task_consume_send_message_event_fn" {
  function_name = "offline_message_notifier_task_consume_send_message_event_fn"
  role          = aws_iam_role.offline_message_notifier_task_consume_send_message_event_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/task_consume_send_message_event"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
    }
  }
}
### End consume event state handler ###

### Start Reset Offline Message Handler ###
resource "aws_lambda_function" "offline_message_notifier_task_reset_offline_message_fn" {
  function_name = "offline_message_notifier_task_reset_offline_message_fn"
  role          = aws_iam_role.offline_message_notifier_task_reset_offline_message_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/task_reset_offline_message"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
    }
  }
}
### End Reset Offline Message Handler ###

### Start Record Offline Message Handler ###
resource "aws_lambda_function" "offline_message_notifier_task_record_offline_message_fn" {
  function_name = "offline_message_notifier_task_record_offline_message_fn"
  role          = aws_iam_role.offline_message_notifier_task_record_offline_message_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/task_record_offline_message"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION = data.aws_region.current.name
    }
  }
}
### End Record Offline Message Handler ###

### Start Trigger Email Handler ###
resource "aws_lambda_function" "offline_message_notifier_task_trigger_offline_email_fn" {
  function_name = "offline_message_notifier_task_trigger_offline_email_fn"
  role          = aws_iam_role.offline_message_notifier_task_trigger_offline_email_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/task_trigger_offline_email"
    ]
  }

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      API_REGION               = data.aws_region.current.name
      "EMAIL_TEMPLATE"         = var.email_template_new_message_open
      "CONFIGURATION_SET_NAME" = var.configuration_set_name
      "EMAIL_IDENTITY"         = var.email_identity
    }
  }
}
### End Trigger Email Handler ###

resource "aws_lambda_function" "sonar_service_support_chat_session_star" {
  function_name = "sonar_service_support_chat_session_star"
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"
  role          = aws_iam_role.sonar_service_support_chat_session_star.arn
  timeout       = 29

  image_config {
    entry_point = ["/lambda/chat_session_star"]
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

# start loop online internal users #
resource "aws_lambda_function" "sonar_service_support_loop_online_internal_users" {
  function_name = "sonar_service_support_loop_online_internal_users"
  role          = aws_iam_role.sonar_service_support_loop_online_internal_users.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = [
      "/lambda/loop_online_internal_users"
    ]
  }

  tracing_config {
    mode = "Active"
  }
}
# end loop online internal users #