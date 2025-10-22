// todo: 途中

# 実行ロール（Lambda 用）
data "aws_iam_policy_document" "assume_lambda" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals { type = "Service", identifiers = ["lambda.amazonaws.com"] }
  }
}

resource "aws_iam_role" "this" {
  name               = "${var.project}-${var.name}-role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda.json
  tags               = var.tags
}

# CloudWatch Logs への出力（最小権限）
data "aws_iam_policy_document" "logs" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "logs" {
  name   = "${var.project}-${var.name}-logs"
  policy = data.aws_iam_policy_document.logs.json
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.this.name
  policy_arn = aws_iam_policy.logs.arn
}

# Bedrock KB Ingestion の起動（必要に応じて絞り込み可）
data "aws_iam_policy_document" "bedrock" {
  statement {
    effect = "Allow"
    actions = [
      "bedrock:StartIngestionJob",
      "bedrock:GetIngestionJob",
      "bedrock:ListIngestionJobs"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "bedrock" {
  name   = "${var.project}-${var.name}-bedrock"
  policy = data.aws_iam_policy_document.bedrock.json
}

resource "aws_iam_role_policy_attachment" "bedrock" {
  role       = aws_iam_role.this.name
  policy_arn = aws_iam_policy.bedrock.arn
}

# Lambda 関数本体
resource "aws_lambda_function" "this" {
  function_name = "${var.project}-${var.name}"
  role          = aws_iam_role.this.arn
  runtime       = var.runtime               # 例: "go1.x"
  architectures = var.architectures         # 例: ["arm64"]
  handler       = var.handler               # 例: "bootstrap"
  filename      = var.zip_path              # 例: "../../apps/kb_ingestor/build/kb_ingestor.zip"
  source_code_hash = filebase64sha256(var.zip_path)

  timeout     = 30
  memory_size = 256

  environment {
    variables = {
      AWS_REGION        = var.region
      KNOWLEDGE_BASE_ID = var.knowledge_base_id
      DATA_SOURCE_ID    = var.data_source_id
    }
  }
}

# （任意）同時実行ガードなどを入れたい場合に備えて出力
resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${aws_lambda_function.this.function_name}"
  retention_in_days = var.log_retention_days
  tags              = var.tags
}
