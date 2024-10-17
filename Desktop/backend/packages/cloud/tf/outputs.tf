output "lambda_file_upload_loop" {
  value = aws_lambda_function.sonar_service_cloud_file_upload_loop
}

output "lambda_file_download_loop" {
  value = aws_lambda_function.sonar_service_cloud_file_download_loop
}

output "lambda_file_upload_loop_v2" {
  value = aws_lambda_function.sonar_service_cloud_file_upload_loop_v2
}

output "lambda_file_download_loop_v2" {
  value = aws_lambda_function.sonar_service_cloud_file_download_loop_v2
}

output "lambda_file_download_web" {
  value = aws_lambda_function.sonar_service_cloud_file_download_web
}

output "lambda_file_upload_web" {
  value = aws_lambda_function.sonar_service_cloud_file_upload_web
}

output "lambda_get_file" {
  value = aws_lambda_function.sonar_service_cloud_get_file
}

output "lambda_associate_file" {
  value = var.live_env ? aws_lambda_function.sonar_service_cloud_associate_file_live_env[0] : aws_lambda_function.sonar_service_cloud_associate_file_dev_env[0]
}

output "lambda_delete_file" {
  value = aws_lambda_function.sonar_service_cloud_delete_file
}

output "lambda_pre_signed_upload_url_web" {
  value = aws_lambda_function.sonar_service_cloud_pre_signed_upload_url_web
}

output "lambda_pre_signed_upload_url_loop" {
  value = aws_lambda_function.sonar_service_cloud_pre_signed_upload_url_loop
}
