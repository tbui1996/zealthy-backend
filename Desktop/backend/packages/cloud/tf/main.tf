terraform {
  required_version = ">= 0.15.0"
}

locals {
  live_envs_web_download = tomap({
    dev  = "https://api.circulo.dev"
    test = "https://api.test.circulosonar.com"
    prod = "https://api.circulosonar.com"
  })
  live_envs_loop_download = tomap({
    dev  = "https://loop-api.circulo.dev"
    test = "https://loop-api.test.circulosonar.com"
    prod = "https://loop-api.circulosonar.com"
  })
  live_envs_upload = tomap({
    dev  = "https://circulo.dev"
    test = "https://test.circulosonar.com"
    prod = "https://circulosonar.com"
  })
}

#### common start
data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

data "aws_iam_policy_document" "lambda_execution_document" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }
}

data "aws_iam_policy_document" "apigw_execution_document" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["apigateway.amazonaws.com"]
      type        = "Service"
    }
  }
}
#### common end
