#!/bin/bash
terraform import module.ecr_router_lambda.aws_ecr_repository.repo circulo/sonar-backend-router-lambda
terraform import module.ecr_forms_lambda.aws_ecr_repository.repo circulo/sonar-backend-forms-lambda
terraform import module.ecr_search_lambda.aws_ecr_repository.repo circulo/sonar-backend-search-lambda
terraform import module.ecr_pearls_lambda.aws_ecr_repository.repo circulo/sonar-backend-pearls-lambda
terraform import module.ecr_support_lambda.aws_ecr_repository.repo circulo/sonar-backend-support-lambda
terraform import module.ecr_users_lambda.aws_ecr_repository.repo circulo/sonar-backend-users-lambda