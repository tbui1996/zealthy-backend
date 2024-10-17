resource "aws_security_group" "external_service_security_group" {
  name   = "sonar_external_service_security_group"
  vpc_id = module.vpc.vpc_id

  description = "Security group for allowing the VPC to access external resources (i.e s3, sqs)"

  ingress {
    description = "Allow VPC to reach external service"
    from_port   = 443
    to_port     = 443

    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    description = "Allow for response back to VPC"
    from_port   = 443
    to_port     = 443

    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "Sonar External Resources Security Group"
    Environment = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }
}
