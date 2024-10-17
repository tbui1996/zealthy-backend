module "offline_message_notifier" {
  source = "./modules/offline-message-notifier"

  task_consume_send_message_event = module.support.task_consume_send_message_event
  task_is_external_user_online    = module.router.task_is_external_user_online
  task_is_internal_user_online    = module.global.task_is_internal_user_online
  task_reset_offline_message      = module.support.task_reset_offline_message
  task_record_offline_message     = module.support.task_record_offline_message
  task_trigger_offline_email      = module.support.task_trigger_offline_email
  task_get_external_user          = module.users.task_get_external_user
  task_get_internal_user          = module.users.task_get_internal_user

  event_bus_name = aws_cloudwatch_event_bus.service_events.name
}