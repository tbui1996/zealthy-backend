terraform {
  required_version = ">= 0.15.0"
}

data "aws_region" "current" {}

resource "aws_acm_certificate" "sonar_cert" {
  domain_name       = "${var.domain_prefix}.${var.domain_name}"
  validation_method = "DNS"

  tags = {
    Name          = "${var.domain_prefix}.${var.domain_name}"
    ProductDomain = var.product_domain
    Environment   = var.environment
    Description   = "Certificate for ${var.domain_prefix}.${var.domain_name}"
    ManagedBy     = "terraform"
  }
}

resource "aws_acm_certificate_validation" "validate_http" {
  certificate_arn         = aws_acm_certificate.sonar_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.http_route53 : record.fqdn]
}

resource "aws_route53_record" "http_route53" {
  for_each = {
    for dvo in aws_acm_certificate.sonar_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
  name            = each.value.name
  records         = [each.value.record]
  type            = each.value.type
  zone_id         = var.hosted_zone_id
  ttl             = 60
  allow_overwrite = true
}

resource "aws_apigatewayv2_domain_name" "http_domain" {
  domain_name = "${var.domain_prefix}.${var.domain_name}"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.sonar_cert.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "http_mapping" {
  api_id      = var.api_id
  domain_name = aws_apigatewayv2_domain_name.http_domain.id
  stage       = var.stage_id
}

resource "aws_route53_record" "http_regional_record" {
  name           = aws_apigatewayv2_domain_name.http_domain.domain_name
  type           = "A"
  zone_id        = var.hosted_zone_id
  set_identifier = var.set_identifier

  alias {
    name                   = aws_apigatewayv2_domain_name.http_domain.domain_name_configuration.0.target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.http_domain.domain_name_configuration.0.hosted_zone_id
    evaluate_target_health = false
  }

  latency_routing_policy {
    region = data.aws_region.current.name
  }
}

