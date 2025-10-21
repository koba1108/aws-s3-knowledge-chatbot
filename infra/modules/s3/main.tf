resource "aws_s3_bucket" "main" {
  bucket = var.knowledge_bucket_name
}

resource "aws_s3_bucket_versioning" "main" {
  bucket = aws_s3_bucket.main.id
  versioning_configuration { status = "Enabled" }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "main" {
  bucket = aws_s3_bucket.main.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "main" {
  bucket                  = aws_s3_bucket.main.id
  block_public_acls       = true
  block_public_policy     = true
  restrict_public_buckets = true
  ignore_public_acls      = true
}

resource "aws_s3_bucket" "log_destination" {
  bucket        = var.logging_destination_bucket
  force_destroy = false
}

resource "aws_s3_bucket_server_side_encryption_configuration" "log_destination_encryption" {
  bucket = aws_s3_bucket.log_destination.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "log_destination_block" {
  bucket                  = aws_s3_bucket.log_destination.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "log_destination_versioning" {
  bucket = aws_s3_bucket.log_destination.id

  versioning_configuration {
    status = "Suspended" # ログ保存用途なら versioning を使わないことも多い
  }
}
