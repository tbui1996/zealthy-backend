variable "task_consume_send_message_event" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_is_external_user_online" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_is_internal_user_online" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_reset_offline_message" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_record_offline_message" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_trigger_offline_email" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_get_external_user" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "task_get_internal_user" {
  type = object({
    arn           = string,
    function_name = string,
  })
}

variable "event_bus_name" {
  type = string
}