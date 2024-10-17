// Configures SES for the domain passed in through environment variables

resource "aws_ses_domain_identity" "default" {
  domain = var.domain_name
}

// DEVELOPMENT ONLY
resource "aws_ses_email_identity" "developer" {
  count = local.live_env || var.okta_username == null ? 0 : 1
  email = var.okta_username
}

resource "aws_route53_record" "ses_verification_record" {
  zone_id = var.hosted_zone_id
  name    = "_amazonses.${aws_ses_domain_identity.default.id}"
  type    = "TXT"
  ttl     = "600"
  records = [aws_ses_domain_identity.default.verification_token]
}

resource "aws_ses_domain_identity_verification" "default" {
  domain     = aws_ses_domain_identity.default.id
  depends_on = [aws_route53_record.ses_verification_record]
}

resource "aws_ses_domain_dkim" "default" {
  domain = aws_ses_domain_identity.default.domain
}

resource "aws_route53_record" "ses_dkim_record" {
  count   = 3
  zone_id = var.hosted_zone_id
  name    = "${element(aws_ses_domain_dkim.default.dkim_tokens, count.index)}._domainkey"
  type    = "CNAME"
  ttl     = "600"
  records = ["${element(aws_ses_domain_dkim.default.dkim_tokens, count.index)}.dkim.amazonses.com"]
}

resource "aws_ses_configuration_set" "default_event_publishing" {
  name                       = "sonar_event_publishing"
  reputation_metrics_enabled = true
}

resource "aws_ses_event_destination" "cloudwatch" {
  name = "sonar_event_publishing_destination"

  configuration_set_name = aws_ses_configuration_set.default_event_publishing.name
  enabled                = true
  matching_types         = ["send", "reject", "bounce", "complaint", "delivery", "open", "click"]

  // include a "type" in the message tags. For example, "circulator_unread_messages"
  cloudwatch_destination {
    default_value  = "default"
    dimension_name = "type"
    value_source   = "messageTag"
  }
}