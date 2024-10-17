output "receive_queue_kms_arn" {
  description = "ARN for the receive queue kms"
  value       = aws_kms_key.service_receive_kms_key.arn
}

output "receive_queue_arn" {
  description = "ARN for the receive queue"
  value       = aws_sqs_queue.service_receive.arn
}

output "send_queue_kms_arn" {
  description = "ARN for the send queue kms"
  value       = aws_kms_key.service_send_kms_key.arn
}

output "send_queue_arn" {
  description = "ARN for the send queue"
  value       = aws_sqs_queue.service_send.arn
}
