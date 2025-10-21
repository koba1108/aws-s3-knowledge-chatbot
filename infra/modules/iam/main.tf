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