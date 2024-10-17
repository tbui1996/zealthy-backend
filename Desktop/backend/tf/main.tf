terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.63.0"
    }
    okta = {
      source  = "okta/okta"
      version = "3.13.3"
    }
    # postgresql = {
    #   source  = "cyrilgdn/postgresql"
    #   version = "1.14.0"
    # }
  }

  backend "http" {}
}

### common start
data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
### common end


### Locals
locals {
  live_env = contains(["dev", "test", "prod"], var.environment)
}

### Providers
provider "aws" {
  region = "us-east-2"
}

provider "aws" {
  alias  = "use1"
  region = "us-east-1"
}

### ECR Modules
module "ecr_router_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-router-lambda"
}

module "ecr_forms_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-forms-lambda"
}

module "ecr_support_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-support-lambda"
}

module "ecr_global_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-global-lambda"
}

module "ecr_users_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-users-lambda"
}

module "ecr_cloud_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-cloud-lambda"
}

module "ecr_feature_flags_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-feature_flags-lambda"
}

module "ecr_patient_lambda" {
  source = "./modules/ecr-repository"
  name   = "circulo/sonar-backend-patient-lambda"
}

# AWS Gateway Configuration
resource "aws_api_gateway_account" "global" {
  cloudwatch_role_arn = aws_iam_role.apigw_cloudwatch_global.arn
}

data "aws_iam_policy_document" "api_gateway_cloudwatch_global_assume_document" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "apigw_cloudwatch_global" {
  name = "api_gateway_cloudwatch_global"

  assume_role_policy = data.aws_iam_policy_document.api_gateway_cloudwatch_global_assume_document.json
}

data "aws_iam_policy_document" "api_gateway_cloudwatch_global_policy_document" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents"
    ]
    resources = [
      # Ignoring the aws-iam-no-policy-wildcards rule because this is used for a global api gateway setting
      "*" #tfsec:ignore:AWS099
    ]
  }
}

resource "aws_iam_role_policy" "apigw_cloudwatch_global" {
  name   = "api_gateway_cloudwatch_global_policy"
  role   = aws_iam_role.apigw_cloudwatch_global.id
  policy = data.aws_iam_policy_document.api_gateway_cloudwatch_global_policy_document.json
}
