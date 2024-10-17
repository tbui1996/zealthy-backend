resource "aws_kms_key" "sonar_cloud_s3_bucket_logs_kms_key" {
  description             = "KMS key for sonar-cloud-files-${var.domain_name}-logs S3 bucket"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_s3_bucket" "sonar_cloud_s3_bucket_logs" {
  bucket = "sonar-cloud-files${var.domain_name}-logs"
  acl    = "log-delivery-write"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.sonar_cloud_s3_bucket_logs_kms_key.id
        sse_algorithm     = "aws:kms"
      }
    }
  }

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "sonar_cloud_logs_s3_bucket_logs_public_access_block" {
  bucket                  = aws_s3_bucket.sonar_cloud_s3_bucket_logs.id
  block_public_policy     = true
  block_public_acls       = true
  restrict_public_buckets = true
  ignore_public_acls      = true
}

resource "aws_kms_key" "sonar_cloud_s3_bucket_kms_key" {
  description             = "KMS key for sonar-cloud-files-${var.domain_name} S3 bucket"
  deletion_window_in_days = 10
  enable_key_rotation     = true
}

resource "aws_s3_bucket" "sonar_cloud" {
  bucket = "sonar-cloud-files-${var.domain_name}"
  acl    = "private"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.sonar_cloud_s3_bucket_kms_key.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }

  versioning {
    enabled = true
  }

  logging {
    target_bucket = aws_s3_bucket.sonar_cloud_s3_bucket_logs.id
    target_prefix = "s3/"
  }

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET"]
    allowed_origins = [
      lookup(local.live_envs_web_download, var.environment, "https://api.${var.environment}.circulo.dev"),
      lookup(local.live_envs_loop_download, var.environment, "https://loop-api.${var.environment}.circulo.dev"),
      lookup(local.live_envs_upload, var.environment, "http://localhost:8080"),
    ]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket_public_access_block" "sonar_cloud_s3_bucket_public_access_block" {
  bucket                  = aws_s3_bucket.sonar_cloud.id
  block_public_policy     = true
  block_public_acls       = true
  restrict_public_buckets = true
  ignore_public_acls      = true
}
