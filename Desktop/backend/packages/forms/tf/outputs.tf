/*
output "routes_web" {
  value = [
    jsonencode(module.forms_get.sonar_gateway_route),
    jsonencode(module.forms_send.sonar_gateway_route),
    jsonencode(module.forms_list.sonar_gateway_route),
    jsonencode(module.forms_create.sonar_gateway_route)
  ]
}
*/

output "lambda_count" {
  value = aws_lambda_function.sonar_service_forms_count
}

output "lambda_create" {
  value = aws_lambda_function.sonar_service_forms_create
}

output "lambda_get" {
  value = aws_lambda_function.sonar_service_forms_get
}

output "lambda_list" {
  value = aws_lambda_function.sonar_service_forms_list
}

output "lambda_response" {
  value = aws_lambda_function.sonar_service_forms_response
}

output "lambda_send" {
  value = aws_lambda_function.sonar_service_forms_send
}

output "lambda_delete" {
  value = aws_lambda_function.sonar_service_forms_delete
}

output "lambda_edit" {
  value = aws_lambda_function.sonar_service_forms_edit
}

output "lambda_close" {
  value = aws_lambda_function.sonar_service_forms_close
}

output "lambda_receive" {
  value = aws_lambda_function.sonar_service_forms_receive
}