output "s3_bucket_name" {
  description = "S3 bucket name for knowledge base"
  value       = aws_s3_bucket.knowledge_base.id
}

output "s3_bucket_arn" {
  description = "S3 bucket ARN for knowledge base"
  value       = aws_s3_bucket.knowledge_base.arn
}

output "knowledge_base_id" {
  description = "Bedrock Knowledge Base ID"
  value       = aws_bedrockagent_knowledge_base.main.id
}

output "knowledge_base_arn" {
  description = "Bedrock Knowledge Base ARN"
  value       = aws_bedrockagent_knowledge_base.main.arn
}

output "data_source_id" {
  description = "Bedrock Knowledge Base Data Source ID"
  value       = aws_bedrockagent_data_source.main.id
}

output "backend_role_arn" {
  description = "IAM role ARN for backend service"
  value       = aws_iam_role.backend_role.arn
}

output "opensearch_collection_endpoint" {
  description = "OpenSearch Serverless collection endpoint"
  value       = aws_opensearchserverless_collection.kb_collection.collection_endpoint
}
