resource "aws_cloudwatch_log_group" "monitoring" {
  name              = var.log_group_name
  retention_in_days = 30
  tags = {
    Project = var.project
  }
}

resource "aws_s3_bucket_logging" "access_logs" {
  for_each      = { for b in var.s3_buckets_to_log : b => b }
  bucket        = each.value
  target_bucket = var.logging_destination_bucket
  target_prefix = "${each.value}/access-logs/"
}
