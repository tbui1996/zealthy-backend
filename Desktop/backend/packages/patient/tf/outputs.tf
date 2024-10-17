output "lambda_list_appointments" {
  value = aws_lambda_function.sonar_service_patient_list_appointments
}

output "lambda_create_appointments" {
  value = aws_lambda_function.sonar_service_patient_create_appointments
}
output "lambda_edit_appointments" {
  value = aws_lambda_function.sonar_service_patient_edit_appointments
}

output "lambda_delete_appointments" {
  value = aws_lambda_function.sonar_service_patient_delete_appointments
}

output "lambda_list_patients" {
  value = aws_lambda_function.sonar_service_patient_list_patients
}

output "lambda_create_patients" {
  value = aws_lambda_function.sonar_service_patient_create_patients
}
output "lambda_patch_patients" {
  value = aws_lambda_function.sonar_service_patient_patch_patients
}


output "lambda_list_agency_providers" {
  value = aws_lambda_function.sonar_service_patient_list_agency_providers
}

output "lambda_create_agency_providers" {
  value = aws_lambda_function.sonar_service_patient_create_agency_providers
}

output "lambda_edit_agency_providers" {
  value = aws_lambda_function.sonar_service_patient_edit_agency_providers
}