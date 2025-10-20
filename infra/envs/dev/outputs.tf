output "main_knowledge_bucket_name" {
  value = module.s3_knowledge.bucket_name
}

output "main_knowledge_bucket_arn" {
  value = module.s3_knowledge.bucket_arn
}

# debug用
output "current_account" { value = data.aws_caller_identity.me.account_id }
output "current_arn"     { value = data.aws_caller_identity.me.arn }
