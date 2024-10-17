resource "aws_sfn_state_machine" "offline_message_notifier_sfn" {
  name     = "offline_message_notifier_state_machine"
  role_arn = aws_iam_role.offline_message_notifier_sfn_role.arn
  type     = "EXPRESS"
  definition = jsonencode({
    Comment : "This state machine will handle control logic for sending offline message notifications",
    StartAt : "Reducer",
    States : {
      Reducer : {
        Type : "Choice",
        Choices : [
          {
            Variable : "$['detail-type']",
            StringEquals : "message_sent_event",
            Next : "ConsumeEvent"
          },
          {
            Variable : "$['detail-type']",
            StringEquals : "connection_created_event",
            Next : "ResetOfflineMessage"
          }
        ],
        Default : "UnknownEventType",
        Comment : "This step function is triggered by two separate events, this state will look at the detail-type of the event and decide which state to go to next."
      },
      ResetOfflineMessage : {
        Type : "Task",
        Resource : "${var.task_reset_offline_message.arn}",
        InputPath : "$.detail.UserID",
        Next : "Done",
        Comment : "If a connection was created for a user, any previous OfflineMessageNotification records for the user should be removed."
      },
      ConsumeEvent : {
        Type : "Task",
        Resource : "${var.task_consume_send_message_event.arn}",
        Next : "IsUserOnlineServiceChoice",
        Comment : "This handler validates and parses the event for the rest of the flow to use."
      },
      IsUserOnlineServiceChoice : {
        Type : "Choice",
        Choices : [{
          Variable : "$$.Execution.Input.detail.SenderType",
          StringEquals : "internal",
          Next : "IsExternalUserOnline"
          }, {
          Variable : "$$.Execution.Input.detail.SenderType",
          StringEquals : "loop",
          Next : "IsInternalUserOnline"
        }]
        Default : "Done",
        Comment : "This choice decides which handler should be used to check if the receiver of the message is online."
      },
      IsExternalUserOnline : {
        Type : "Task",
        Resource : "${var.task_is_external_user_online.arn}",
        Next : "IsUserOnlineChoice",
        InputPath : "$.ReceiverId",
        Comment : "Check whether or not an external user is online."
      },
      IsInternalUserOnline : {
        Type : "Task",
        Resource : "${var.task_is_internal_user_online.arn}",
        Next : "IsUserOnlineChoice",
        InputPath : "$.ReceiverId",
        Comment : "Check whether or not an internal user is online."
      },
      IsUserOnlineChoice : {
        Type : "Choice",
        Choices : [
          {
            Variable : "$.IsOnline",
            BooleanEquals : true,
            Next : "Done"
          },
        ],
        Default : "RecordOfflineMessage",
        Comment : "If the user is offline, we should record an OfflineMessageNotification record to track that we should notifiy them. Otherwise, we can remove previous notifications if they exist."
      },
      RecordOfflineMessage : {
        Type : "Task",
        Resource : "${var.task_record_offline_message.arn}",
        Next : "ShouldTriggerEmailChoice",
        InputPath : "$.UserId",
        Comment : "Track that a user should be notified. If a notification was already scheduled, this returns false for Created."
      },
      ShouldTriggerEmailChoice : {
        Type : "Choice",
        Choices : [
          {
            Variable : "$.Created",
            BooleanEquals : true,
            Next : "TriggerEmailWait"
          },
        ],
        Default : "Done",
        Comment : "If a notification is already scheduled, do not trigger an email, otherwise trigger an email."
      },
      TriggerEmailWait : {
        Type : "Wait",
        Seconds : 300,
        Next : "UserLookupChoice",
        Comment : "Delay the email by 5 minutes"
      },
      UserLookupChoice : {
        Type : "Choice",
        Choices : [{
          Variable : "$$.Execution.Input.detail.SenderType",
          StringEquals : "internal",
          Next : "GetExternalUserInfo"
          }, {
          Variable : "$$.Execution.Input.detail.SenderType",
          StringEquals : "loop",
          Next : "GetInternalUserInfo"
        }]
        Default : "Done",
        Comment : "Decide which handler to use to lookup user info. This state is using context ($$) rather than input ($)"
      },
      GetExternalUserInfo : {
        Type : "Task",
        Resource : "${var.task_get_external_user.arn}",
        Next : "TriggerEmail",
        InputPath : "$$.Execution.Input.detail.ReceiverId",
        Comment : "Get user information for external user"
      },
      GetInternalUserInfo : {
        Type : "Task",
        Resource : "${var.task_get_internal_user.arn}",
        Next : "TriggerEmail",
        InputPath : "$$.Execution.Input.detail.ReceiverId",
        Comment : "Get user information for internal user"
      },
      TriggerEmail : {
        Type : "Task",
        Resource : "${var.task_trigger_offline_email.arn}",
        Next : "Done",
        Comment : "Send the user an email and mark their notification as sent"
      },
      UnknownEventType : {
        Type : "Fail",
        Cause : "The event detail type was unknown.",
      }
      Done : {
        Type : "Succeed"
      }
    }
  })

  logging_configuration {
    log_destination        = "${aws_cloudwatch_log_group.offline_message_notifier_sfn_logs.arn}:*"
    level                  = "ALL"
    include_execution_data = true
  }
}