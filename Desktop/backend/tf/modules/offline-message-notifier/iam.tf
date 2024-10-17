# Start Step Function Permissions #
resource "aws_iam_role" "offline_message_notifier_sfn_role" {
  name               = "offline_message_notifier_sfn_role"
  assume_role_policy = data.aws_iam_policy_document.sfn_execution_document.json
}

data "aws_iam_policy_document" "offline_message_notifier_sfn_document" {
  statement {
    actions = ["lambda:InvokeFunction"]

    resources = [
      // support
      var.task_consume_send_message_event.arn,
      var.task_reset_offline_message.arn,
      var.task_record_offline_message.arn,
      var.task_trigger_offline_email.arn,
      // router
      var.task_is_external_user_online.arn,
      var.task_is_internal_user_online.arn,
      // users
      var.task_get_external_user.arn,
      var.task_get_internal_user.arn
    ]
  }

  statement {
    actions = [
      "logs:CreateLogDelivery",
      "logs:GetLogDelivery",
      "logs:UpdateLogDelivery",
      "logs:DeleteLogDelivery",
      "logs:ListLogDeliveries",
      "logs:PutResourcePolicy",
      "logs:DescribeResourcePolicies",
      "logs:DescribeLogGroups"
    ]
    # I spent a ton of time trying to get this working, between trying to enable KMS and getting this resource right I ended up giving up since it's just logs
    # open to ideas!
    #tfsec:ignore:AWS099
    resources = [
      "*"
    ]
  }
}

resource "aws_iam_role_policy" "offline_message_notifier_sfn_policy" {
  name = "offline_message_notifier_sfn_policy-policy"
  role = aws_iam_role.offline_message_notifier_sfn_role.id

  policy = data.aws_iam_policy_document.offline_message_notifier_sfn_document.json
}

# End Step Function Permissions #

# Start Event Rule Permissions #
resource "aws_iam_role" "send_message_event_rule_role" {
  name               = "send_message_event_rule_role"
  assume_role_policy = data.aws_iam_policy_document.events_execution_document.json
}

data "aws_iam_policy_document" "invoke_offline_message_notifier_policy_document" {
  statement {
    actions = ["states:StartExecution"]

    resources = [aws_sfn_state_machine.offline_message_notifier_sfn.arn]
  }
}

resource "aws_iam_role_policy" "invoke_offline_message_notifier_policy" {
  name = "invoke_offline_message_notifier_policy"
  role = aws_iam_role.send_message_event_rule_role.id

  policy = data.aws_iam_policy_document.invoke_offline_message_notifier_policy_document.json
}
# End Event Rule Permissions #