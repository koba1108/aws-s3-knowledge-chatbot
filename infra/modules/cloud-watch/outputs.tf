output "log_group_arn" {
  value       = aws_cloudwatch_log_group.monitoring.arn
  description = "ARN of the CloudWatch Log Group"
}
