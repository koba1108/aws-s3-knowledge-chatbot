output "bucket_name" { value = aws_s3_bucket.main.bucket }
output "bucket_arn" { value = aws_s3_bucket.main.arn }
output "logging_bucket_name" { value = var.logging_destination_bucket }
