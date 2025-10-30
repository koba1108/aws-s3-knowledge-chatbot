data "aws_caller_identity" "me" {}

locals {
  caller_arn      = data.aws_caller_identity.me.arn
  caller_parts    = split("/", local.caller_arn)
  caller_username = local.caller_parts[length(local.caller_parts) - 1]
}

data "aws_iam_policy_document" "kb_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["bedrock.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "kb_role" {
  name               = "${var.project}-kb-role"
  assume_role_policy = data.aws_iam_policy_document.kb_assume.json
}

data "aws_iam_policy_document" "kb_access" {
  statement {
    sid       = "S3Read"
    actions   = ["s3:ListBucket"]
    resources = [var.bucket_arn]
  }
  statement {
    sid       = "S3Get"
    actions   = ["s3:GetObject"]
    resources = ["${var.bucket_arn}/*"]
  }
  statement {
    sid       = "AOSSApiAccess"
    actions   = ["aoss:APIAccessAll"]
    resources = ["*"] # API権限はリソース指定不可のためワイルドカード
  }
  statement {
    sid = "BedrockInvokeEmbedding"
    actions = ["bedrock:InvokeModel"]
    resources = [
      "arn:aws:bedrock:ap-northeast-1::foundation-model/amazon.titan-embed-text-v2:0"
    ]
  }
}

resource "aws_iam_policy" "kb_access" {
  name   = "${var.project}-kb-access"
  policy = data.aws_iam_policy_document.kb_access.json
}

resource "aws_iam_role_policy_attachment" "attach" {
  role       = aws_iam_role.kb_role.name
  policy_arn = aws_iam_policy.kb_access.arn
}


resource "aws_iam_policy" "oss_dashboard_access" {
  name        = "opensearch-serverless-dashboard-access"
  description = "Allow access to OpenSearch Serverless Dashboards and API"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["aoss:APIAccessAll"]
        Resource = "arn:aws:aoss:${var.region}:${var.current_account}:collection/${var.oss_collection_id}"
      },
      # ダッシュボードUIへのアクセス
      {
        Effect   = "Allow"
        Action   = ["aoss:DashboardsAccessAll"]
        Resource = "*"
      }
    ]
  })
}

# --- RAG 実行用ポリシー（実行ユーザー/ロール向け） ---
data "aws_iam_policy_document" "rag_runtime" {
  statement {
    sid     = "KBReadRAG"
    effect  = "Allow"
    actions = [
      "bedrock:RetrieveAndGenerate",
      "bedrock:RetrieveAndGenerateStream",
    ]
    resources = ["*"] # 後で KB ARN に絞るとより安全
  }

  statement {
    sid     = "InvokeGenModel"
    effect  = "Allow"
    actions = [
      "bedrock:InvokeModel",
      "bedrock:InvokeModelWithResponseStream",
    ]
    resources = ["*"] # 後で foundation-model ARN に絞るとより安全
  }
}

resource "aws_iam_policy" "rag_runtime" {
  name   = "${var.project}-rag-runtime"
  policy = data.aws_iam_policy_document.rag_runtime.json
}

# Terraform 実行中の IAM ユーザー（caller_username）に付与
resource "aws_iam_user_policy_attachment" "attach_rag_to_caller" {
  user       = local.caller_username
  policy_arn = aws_iam_policy.rag_runtime.arn
}