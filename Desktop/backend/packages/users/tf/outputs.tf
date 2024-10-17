output "external_cognito_client_id" {
  description = "External Cognito Client ID"
  value       = aws_cognito_user_pool_client.externals.id
}

output "aws_cognito_user_pool_externals_jwks" {
  description = "Location of the JSON Web Key Set (JWKS) to be used to verify token signatures from external cognito user pool"
  value       = "https://${aws_cognito_user_pool.externals.endpoint}/.well-known/jwks.json"
}

output "aws_cognito_user_pool_internals_jwks" {
  description = "Location of the JSON Web Key Set (JWKS) to be used to verify token signatures from internal cognito user pool"
  value       = "https://${aws_cognito_user_pool.internals.endpoint}/.well-known/jwks.json"
}

output "lambda_external_sign_in" {
  value = aws_lambda_function.sonar_service_users_external_sign_in
}

output "lambda_external_sign_up" {
  value = aws_lambda_function.sonar_service_users_external_sign_up
}

output "lambda_external_refresh" {
  value = aws_lambda_function.sonar_service_users_external_refresh
}

output "lambda_user_list" {
  value = aws_lambda_function.sonar_service_users_user_list
}

output "lambda_revoke_access" {
  value = aws_lambda_function.sonar_service_users_revoke_access
}

output "lambda_receive" {
  value = aws_lambda_function.sonar_service_users_receive
}

output "lambda_get_organizations" {
  value = aws_lambda_function.sonar_service_users_get_organizations
}

output "lambda_update_user" {
  value = aws_lambda_function.sonar_service_update_user
}

output "task_get_external_user" {
  value = aws_lambda_function.sonar_service_users_task_get_external_user
}

output "task_get_internal_user" {
  value = aws_lambda_function.sonar_service_users_task_get_internal_user
}

output "lambda_create_organizations" {
  value = aws_lambda_function.sonar_service_users_create_organizations
}