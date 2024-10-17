output "lambda_chat_session" {
  value = aws_lambda_function.sonar_service_support_chat_session_create
}

output "lambda_pending_chat_session_create" {
  value = aws_lambda_function.sonar_service_support_pending_chat_session_create
}

output "lambda_pending_chat_sessions_get" {
  value = aws_lambda_function.sonar_service_support_pending_chat_sessions_get
}

output "lambda_assign_pending_chat_session" {
  value = aws_lambda_function.sonar_service_support_assign_pending_chat_session
}

output "lambda_chat_messages_get" {
  value = aws_lambda_function.sonar_service_support_chat_messages_get
}

output "lambda_chat_sessions_get" {
  value = aws_lambda_function.sonar_service_support_chat_sessions_get
}

output "lambda_chat_session_update_open" {
  value = aws_lambda_function.sonar_service_support_chat_session_update_open
}

output "lambda_send" {
  value = aws_lambda_function.sonar_service_support_send
}

output "lambda_loop_chat_messages_get" {
  value = aws_lambda_function.sonar_service_support_loop_chat_messages_get
}

output "lambda_loop_chat_sessions_get" {
  value = aws_lambda_function.sonar_service_support_loop_chat_sessions_get
}

output "lambda_loop_send" {
  value = aws_lambda_function.sonar_service_support_loop_send
}

output "lambda_receive" {
  value = aws_lambda_function.sonar_service_support_receive
}

output "lambda_update_chat_notes" {
  value = aws_lambda_function.sonar_service_support_update_chat_notes
}

output "lambda_submit_feedback" {
  value = aws_lambda_function.sonar_service_support_submit_feedback
}

output "task_consume_send_message_event" {
  value = aws_lambda_function.offline_message_notifier_task_consume_send_message_event_fn
}

output "task_reset_offline_message" {
  value = aws_lambda_function.offline_message_notifier_task_reset_offline_message_fn
}

output "task_record_offline_message" {
  value = aws_lambda_function.offline_message_notifier_task_record_offline_message_fn
}

output "task_trigger_offline_email" {
  value = aws_lambda_function.offline_message_notifier_task_trigger_offline_email_fn
}

output "lambda_chat_session_star" {
  value = aws_lambda_function.sonar_service_support_chat_session_star
}

output "lambda_loop_online_internal_users" {
  value = aws_lambda_function.sonar_service_support_loop_online_internal_users
}

output "lambda_loop_patients_get" {
  value = aws_lambda_function.sonar_service_support_loop_patients_get
}