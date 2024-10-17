# I burned a lot of time trying to get logging for the step function and KMS workign well together, I kept getting intermittent TF apply errors
# Since logs are encrypted no matter what, and the KMS key mostly just adds an extra layer of audit, Tony and I discussed not adding a key here
#tfsec:ignore:AWS089
resource "aws_cloudwatch_log_group" "offline_message_notifier_sfn_logs" {
  name = "/aws/states/offline_message_notifier"
}