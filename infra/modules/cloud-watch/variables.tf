variable "project" {
  type        = string
  description = "Project prefix for tagging and naming"
}

variable "log_group_name" {
  type        = string
  description = "Name of the CloudWatch Logs group"
}

variable "s3_buckets_to_log" {
  type        = list(string)
  description = "List of S3 bucket names to enable server access logging"
}

variable "logging_destination_bucket" {
  type        = string
  description = "Bucket name where S3 access logs will be saved (for all other buckets)"
}
