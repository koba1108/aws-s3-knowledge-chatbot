variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "aws-s3-kb-chatbot"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "embedding_model" {
  description = "Bedrock embedding model"
  type        = string
  default     = "amazon.titan-embed-text-v1"
}
