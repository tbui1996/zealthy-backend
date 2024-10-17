output "lambda_create_flag" {
  value = aws_lambda_function.sonar_service_feature_flags_create_flag
}

output "lambda_patch_flag" {
  value = aws_lambda_function.sonar_service_feature_flags_patch_flag
}

output "lambda_list_flags" {
  value = aws_lambda_function.sonar_service_feature_flags_list_flags
}

output "lambda_evaluate" {
  value = aws_lambda_function.sonar_service_feature_flags_evaluate
}

output "lambda_loop_evaluate" {
  value = aws_lambda_function.sonar_service_feature_flags_loop_evaluate
}

output "lambda_delete_flag" {
  value = aws_lambda_function.sonar_service_feature_flags_delete_flag
}
