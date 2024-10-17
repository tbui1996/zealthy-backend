resource "aws_lambda_function" "sonar_service_SERVICE_NAME_LAMBDA_NAME" {
  function_name = "sonar_service_SERVICE_NAME_LAMBDA_NAME"
  role          = aws_iam_role.sonar_service_SERVICE_NAME_LAMBDA_NAME_role.arn
  package_type  = "Image"
  image_uri     = "${var.lambda_image_uri}:${var.image_version}"

  image_config {
    entry_point = ["/lambda/LAMBDA_NAME"]
  }

  tracing_config {
    mode = "Active"
  }
}