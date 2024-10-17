# cleanup user_list_connected
terraform destroy  \
  --target=module.router_user_list_connected.aws_lambda_permission.apigw_http_lambda \
  --target=module.router_user_list_connected.aws_iam_role_policy.api_gateway_invocation_policy \
  --target=module.router_user_list_connected.aws_apigatewayv2_integration.http_lambda_integration \
  --target=module.router_user_list_connected.aws_iam_role.api_gateway_invocation_role \
  --target=module.router_user_list_connected.aws_apigatewayv2_route.http_lambda \
  --target=module.router.aws_lambda_function.api_users_list_connected \
  --target=module.router.aws_iam_role.api_users_list_connected_role \
  --target=module.router.aws_iam_role_policy.api_users_list_connected_logging \
  --target=module.router.aws_iam_role_policy.api_users_list_connected_dynamodb
