data "aws_caller_identity" "me" {}

module "s3_knowledge" {
  source                = "../../modules/s3"
  knowledge_bucket_name = var.knowledge_bucket_name
}
