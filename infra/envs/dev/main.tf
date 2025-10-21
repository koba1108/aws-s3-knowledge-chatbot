data "aws_caller_identity" "me" {}

module "iam" {
  source = "../../modules/iam"

  project           = var.project
  region            = var.region
  current_account   = data.aws_caller_identity.me.account_id
  bucket_arn        = module.s3_knowledge.bucket_arn
  oss_collection_id = module.opensearch_serverless.collection_id
}

module "s3_knowledge" {
  source = "../../modules/s3"

  knowledge_bucket_name      = var.knowledge_bucket_name
  logging_destination_bucket = var.logging_destination_bucket
}

module "opensearch_serverless" {
  source = "../../modules/opensearch-serverless"

  collection_name = var.opensearch_collection_name
  allow_public    = var.opensearch_allow_public
  principals = [
    data.aws_caller_identity.me.arn,
    module.iam.kb_role_arn,
  ]
}

module "bedrock_knowledge_base" {
  source = "../../modules/bedrock-knowledge-base"

  project            = var.project
  role_arn           = module.iam.kb_role_arn
  oss_collection_arn = module.opensearch_serverless.collection_arn
  bucket_arn         = module.s3_knowledge.bucket_arn
}

module "cloud_watch" {
  source = "../../modules/cloud-watch"

  project                    = var.project
  s3_buckets_to_log          = [module.s3_knowledge.bucket_name]
  log_group_name             = "/aws/bedrock/${module.bedrock_knowledge_base.knowledge_base_id}"
  logging_destination_bucket = module.s3_knowledge.logging_bucket_name

  depends_on = [module.s3_knowledge]
}
