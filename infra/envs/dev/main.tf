data "aws_caller_identity" "me" {}

module "s3_knowledge" {
  source = "../../modules/s3"

  knowledge_bucket_name = var.knowledge_bucket_name
}

module "opensearch_serverless" {
  source = "../../modules/opensearch-serverless"

  collection_name = var.opensearch_collection_name
  allow_public    = var.opensearch_allow_public
  principals      = [data.aws_caller_identity.me.arn]
}
