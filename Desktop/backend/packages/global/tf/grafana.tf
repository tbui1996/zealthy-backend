#tfsec:ignore:aws-elbv2-alb-not-public
resource "aws_lb" "sonar_grafana" {
  name = "sonar-grafana-lb"
  # This load balancer needs to accept internet traffic so it is external,
  # but external load balancers get flagged by tfsec, so we tfsec:ignore it.
  # https://tfsec.dev/docs/aws/elbv2/alb-not-public/#aws/elbv2
  internal           = false
  load_balancer_type = "application"
  subnets            = var.vpc_public_subnets
  security_groups    = [aws_security_group.sonar_grafana_lb.id]
  # Invalid headers can be used to hack web applications.
  # Read: https://www.hacksparrow.com/webdev/security/dangers-of-trusting-http-headers.html
  drop_invalid_header_fields = true
}

# This security group is for allowing the internet to hit our load balancer, and vice versa
resource "aws_security_group" "sonar_grafana_lb" {
  description = "Allow inbound traffic to our Grafana instance"
  vpc_id      = var.vpc_id

  ingress {
    description = "Allow all"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # egress is needed so we can pull the official docker grafana image
  egress {
    description = "Allow all"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Security groups can not be destroyed if they are being used.
  # https://shinglyu.com/web/2020/02/06/update-aws-security-groups-with-terraform.html
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_target_group" "target_group" {
  name = "sonar-grafana-target-group"
  # This technically does not need to be port 80, but HTTP is port 80 by default.
  # If we want to change it something other than 80/HTTP, than a lot more configuration is required.
  # Also, at this point, we are already inside of the VPC from our previous security group.
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = var.vpc_id

  health_check {
    matcher = "200"
    path    = "/api/health"
  }
}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_lb.sonar_grafana.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-FS-2018-06"

  certificate_arn = aws_acm_certificate.grafana.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.target_group.arn
  }
}

resource "aws_acm_certificate" "grafana" {
  domain_name       = "grafana.${var.domain_name}"
  validation_method = "DNS"

  tags = {
    Name          = "grafana.${var.domain_name}"
    ProductDomain = "Sonar Grafana"
    Environment   = var.environment
    Description   = "Certificate for grafana.${var.domain_name}"
    ManagedBy     = "terraform"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "validate_grafana" {
  certificate_arn         = aws_acm_certificate.grafana.arn
  validation_record_fqdns = [for record in aws_route53_record.grafana_route53 : record.fqdn]
}

resource "aws_route53_record" "grafana_route53" {
  # When we create the acm_certificate above, acm generates record names, record values, and record types.
  # We use each of those to create our route53 records
  # https://aws.amazon.com/premiumsupport/knowledge-center/route-53-validate-acm-certificates/
  # https://www.oss-group.co.nz/blog/automated-certificates-aws
  for_each = {
    for dvo in aws_acm_certificate.grafana.domain_validation_options : dvo.domain_name => {
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

resource "aws_route53_record" "sonar_service_grafana_regional" {
  zone_id = var.hosted_zone_id
  name    = "grafana.${var.domain_name}"
  type    = "A"

  alias {
    name    = aws_lb.sonar_grafana.dns_name
    zone_id = aws_lb.sonar_grafana.zone_id
    # This health check is a latency check and may be overkill.
    # Enabling this also drives up our costs, but can do later if needed.
    # We're already checking the health of the target group (200 /api/health)
    # https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/health-checks-how-route-53-chooses-records.html
    evaluate_target_health = false
  }
}

resource "aws_kms_key" "sonar_grafana_logs_kms_key" {
  description             = "KMS key for Sonar Grafana cloudwatch logs"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.grafana_logs_kms_key_document.json
}

data "aws_iam_policy_document" "grafana_logs_kms_key_document" {
  policy_id = "sonar-grafana-logs-kms-key"
  statement {
    sid = "Enable IAM User Permissions"
    principals {
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
      type        = "AWS"
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }

  statement {
    sid = "Enable key to be used with Cloud Watch Logs"
    principals {
      identifiers = ["logs.${data.aws_region.current.name}.amazonaws.com"]
      type        = "Service"
    }
    actions = [
      "kms:Encrypt*",
      "kms:Decrypt*",
      "kms:ReEncryptTo",
      "kms:ReEncryptFrom",
      "kms:GenerateDataKey*",
      "kms:Describe*"
    ]
    resources = ["*"]
  }
}
resource "aws_kms_alias" "sonar_grafana_logs_kms_key_alias" {
  name          = "alias/grafana_logs_kms_key"
  target_key_id = aws_kms_key.sonar_grafana_logs_kms_key.arn
}

resource "aws_cloudwatch_log_group" "grafana_logs" {
  name       = "sonar-grafana"
  kms_key_id = aws_kms_key.sonar_grafana_logs_kms_key.arn
}

resource "aws_iam_role" "ecs_task_role" {
  name               = "sonar_service_grafana_ecs_task_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}

resource "aws_ecs_task_definition" "sonar-service-grafana" {
  family = "sonar-service-grafana"

  container_definitions = jsonencode(
    [
      {
        name : "sonar-service-grafana",
        image : "grafana/grafana",
        essential : true,
        portMappings : [
          {
            containerPort : 3000,
            hostPort : 3000
          }
        ],
        memory : 2048,
        cpu : 1024,
        logConfiguration : {
          logDriver : "awslogs",
          options : {
            awslogs-group : "${aws_cloudwatch_log_group.grafana_logs.name}",
            awslogs-region : "${data.aws_region.current.name}",
            awslogs-stream-prefix : "dashboard"
          }
        }
      }
    ]
  )

  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  memory                   = 2048
  cpu                      = 1024
  execution_role_arn       = aws_iam_role.ecs_task_role.arn

  volume {
    name = "grafana-fs"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.grafana_fs.id
      root_directory     = "/var/lib/grafana"
      transit_encryption = "ENABLED"
    }
  }
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_role_policy" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# This security group allows our load balancer to hit the Grafana instance.
# Without it, the internet hits our load balancer, but the load balancer can't
# route you to the Grafana instance on port 3000, resulting in a "service temp unavailable" error.
resource "aws_security_group" "sonar-service-grafana-ecs" {
  description = "Security group for sonar-service-grafana ECS service to allow traffic to/from load balancer"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Allow all"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [aws_security_group.sonar_grafana_lb.id]
  }

  # We need egress here for the ecs service to pull the docker image
  egress {
    description = "Allow all"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_ecs_service" "sonar-service-grafana" {
  name            = "sonar-service-grafana"
  cluster         = var.ecs_cluster_name
  task_definition = aws_ecs_task_definition.sonar-service-grafana.arn
  launch_type     = "FARGATE"
  desired_count   = 1

  network_configuration {
    subnets          = var.vpc_private_subnets
    assign_public_ip = true
    security_groups  = [aws_security_group.sonar-service-grafana-ecs.id]
  }

  # Here we are creating an internal load balancer, which routes traffic
  # from the port 80 of the target group to port 3000 of the container.
  load_balancer {
    target_group_arn = aws_lb_target_group.target_group.arn
    container_name   = aws_ecs_task_definition.sonar-service-grafana.family
    container_port   = 3000
  }
}

resource "aws_efs_file_system" "grafana_fs" {
  encrypted = true
  tags = {
    Name = "Grafana FS"
  }
}

resource "aws_efs_mount_target" "grafana_fs" {
  file_system_id = aws_efs_file_system.grafana_fs.id
  subnet_id      = var.vpc_private_subnets[0]
}
