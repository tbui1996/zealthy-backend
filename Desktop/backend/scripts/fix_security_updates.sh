#!/bin/bash
terraform destroy -target=aws_security_group.rds_security_group
terraform destroy -target=module.cms.aws_security_group.load_balancer
terraform destroy -target=module.search.aws_elasticsearch_domain.search

# Run task deploy a few times
# You have to make a commit message now
