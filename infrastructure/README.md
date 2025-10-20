# Infrastructure

This directory contains Terraform configuration for AWS infrastructure.

## Resources Created

- **S3 Bucket**: Storage for knowledge base documents
- **Bedrock Knowledge Base**: Vector database for semantic search
- **OpenSearch Serverless**: Vector store for embeddings
- **IAM Roles**: Permissions for Bedrock and backend services

## Setup

1. Initialize Terraform:
```bash
terraform init
```

2. Review the plan:
```bash
terraform plan
```

3. Apply the configuration:
```bash
terraform apply
```

## Configuration

Edit `variables.tf` or create a `terraform.tfvars` file:

```hcl
aws_region     = "us-east-1"
project_name   = "aws-s3-kb-chatbot"
environment    = "dev"
embedding_model = "amazon.titan-embed-text-v1"
```

## Outputs

After applying, Terraform will output:
- S3 bucket name and ARN
- Knowledge Base ID and ARN
- Backend IAM role ARN
- OpenSearch collection endpoint
