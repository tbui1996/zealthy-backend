resource "aws_cognito_user_group" "internal_groups" {
  for_each     = var.internal_group_names
  name         = "internals_${each.key}"
  user_pool_id = aws_cognito_user_pool.internals.id
  description  = "User Group for ${each.key}"
}

resource "aws_cognito_user_group" "external_groups" {
  for_each     = var.external_group_names
  name         = "externals_${each.key}"
  user_pool_id = aws_cognito_user_pool.externals.id
  description  = "User Group for ${each.key}"
}