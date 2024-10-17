### Event Rules Start ###
resource "aws_cloudwatch_event_rule" "message_sent_event_rule" {
  name        = "message_sent_event_rule"
  description = "event bridge rule to capture message_sent_events"

  event_bus_name = var.event_bus_name

  role_arn = aws_iam_role.send_message_event_rule_role.arn

  event_pattern = jsonencode({
    "detail-type" : ["message_sent_event", "connection_created_event"]
  })
}
### Event Rules End ###

### Event Targets Start ###
resource "aws_cloudwatch_event_target" "offline_message_notifier_target" {
  event_bus_name = var.event_bus_name
  target_id      = "offline_message_notifier_target"
  rule           = aws_cloudwatch_event_rule.message_sent_event_rule.name
  arn            = aws_sfn_state_machine.offline_message_notifier_sfn.arn

  role_arn = aws_iam_role.send_message_event_rule_role.arn
}
### Event Targets End ###