output "acm_certificate_cert_arn" {
  value = aws_acm_certificate.sonar_cert.arn
}

output "domain_name_id" {
  value = aws_apigatewayv2_domain_name.http_domain.id
}
