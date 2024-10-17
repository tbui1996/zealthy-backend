data "aws_iam_policy_document" "group_policy" {
  for_each = var.group_names

  statement {
    sid = "sonar_group_${each.key}"

    actions = [
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]

    resources = lookup(var.group_routes, each.key)
  }
}

resource "aws_dynamodb_table_item" "group_policy_item" {
  for_each = data.aws_iam_policy_document.group_policy

  table_name = aws_dynamodb_table.sonar_group_policy.name
  hash_key   = aws_dynamodb_table.sonar_group_policy.hash_key

  item = jsonencode({
    group : { S : "${each.key}" },
    policy : { S : "${replace(each.value.json, "/\\n/", "")}" }
  })
}




