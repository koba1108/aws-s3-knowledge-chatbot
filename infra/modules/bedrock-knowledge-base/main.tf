resource "aws_bedrockagent_knowledge_base" "main" {
  name     = "${var.project}-kb"
  role_arn = var.role_arn

  knowledge_base_configuration {
    type = "VECTOR"
    vector_knowledge_base_configuration {
      embedding_model_arn = "arn:aws:bedrock:ap-northeast-1::foundation-model/amazon.titan-embed-text-v2:0"
    }
  }

  storage_configuration {
    type = "OPENSEARCH_SERVERLESS"

    opensearch_serverless_configuration {
      collection_arn    = var.oss_collection_arn
      vector_index_name = "knowledge-base-index" # todo: indexは手作業で作成する必要あり
      field_mapping {
        text_field     = "text"
        vector_field   = "vector"
        metadata_field = "metadata"
      }
    }
  }
}

resource "aws_bedrockagent_data_source" "s3" {
  knowledge_base_id = aws_bedrockagent_knowledge_base.main.id
  name              = "${var.project}-kb-s3"

  depends_on = [aws_bedrockagent_knowledge_base.main]

  data_source_configuration {
    type = "S3"
    s3_configuration {
      bucket_arn         = var.bucket_arn
      inclusion_prefixes = ["knowledge"]
    }
  }

  vector_ingestion_configuration {
    chunking_configuration {
      chunking_strategy = "FIXED_SIZE"
      fixed_size_chunking_configuration {
        max_tokens         = 500
        overlap_percentage = 50
      }
    }
  }
}
