output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.sonar_cluster.name
}

output "webapp_websocket_api_id" {
  description = "Websocket API Internal ID"
  value       = aws_apigatewayv2_api.webapp_websocket_api.id
}

output "webapp_websocket_api_execution_arn" {
  description = "Websocket API execution Internal ARN"
  value       = aws_apigatewayv2_api.webapp_websocket_api.execution_arn
}