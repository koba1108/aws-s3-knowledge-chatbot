output "collection_name" {
  value       = aws_opensearchserverless_collection.vec.name
  description = "OpenSearch Serverless collection name"
}

output "collection_arn" {
  value       = aws_opensearchserverless_collection.vec.arn
  description = "OpenSearch Serverless collection ARN"
}
