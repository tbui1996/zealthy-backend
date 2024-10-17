#!/bin/bash
terraform destroy -target=module.router.aws_dynamodb_table.pending_messages
terraform destroy -target=module.router.aws_dynamodb_table.websocket_connections

# Run task deploy a few times
# You have to make a commit message now
