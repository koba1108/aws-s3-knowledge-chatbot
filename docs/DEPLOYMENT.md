# デプロイメントガイド

## 前提条件

- AWS アカウント
- AWS CLI v2
- Terraform >= 1.0
- Go >= 1.20
- Node.js >= 18
- AWS Bedrock へのアクセス権限

## AWS Bedrock の有効化

1. AWS マネジメントコンソールにログイン
2. Bedrock サービスに移動
3. 使用するリージョン (例: us-east-1) を選択
4. Model access で以下のモデルを有効化:
   - Anthropic Claude 3 Sonnet
   - Amazon Titan Embeddings G1 - Text

## 手順

### ステップ 1: リポジトリのクローン

```bash
git clone https://github.com/koba1108/aws-s3-knowledge-chatbot.git
cd aws-s3-knowledge-chatbot
```

### ステップ 2: インフラストラクチャのデプロイ

```bash
cd infrastructure

# Terraform 初期化
terraform init

# 変数ファイルの作成 (オプション)
cat > terraform.tfvars <<EOF
aws_region     = "us-east-1"
project_name   = "my-chatbot"
environment    = "prod"
EOF

# プランの確認
terraform plan

# デプロイ実行
terraform apply

# 出力値を記録
terraform output
```

出力例:
```
knowledge_base_id = "ABC123XYZ"
s3_bucket_name = "my-chatbot-kb-prod"
backend_role_arn = "arn:aws:iam::123456789012:role/..."
```

### ステップ 3: ドキュメントのアップロード

```bash
# S3バケット名を環境変数に設定
export S3_BUCKET=$(terraform output -raw s3_bucket_name)

# サンプルドキュメントをアップロード
aws s3 cp docs/sample.pdf s3://$S3_BUCKET/
aws s3 cp docs/faq.txt s3://$S3_BUCKET/
```

### ステップ 4: データソースの同期

```bash
# Knowledge Base IDとData Source IDを取得
export KB_ID=$(terraform output -raw knowledge_base_id)
export DS_ID=$(terraform output -raw data_source_id)

# インジェストジョブを開始
aws bedrock-agent start-ingestion-job \
  --knowledge-base-id $KB_ID \
  --data-source-id $DS_ID \
  --region us-east-1

# ジョブの状態を確認
aws bedrock-agent list-ingestion-jobs \
  --knowledge-base-id $KB_ID \
  --data-source-id $DS_ID \
  --region us-east-1
```

インジェストが完了するまで数分かかります。

### ステップ 5: バックエンドのデプロイ

#### 5.1 ローカル実行

```bash
cd ../backend

# 環境変数を設定
export KNOWLEDGE_BASE_ID=$KB_ID
export AWS_REGION=us-east-1
export MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0
export PORT=8080

# 依存関係のインストール
go mod download

# ビルド
go build -o server main.go

# 実行
./server
```

#### 5.2 AWS Lambda へのデプロイ

```bash
# ビルド (Linux用)
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# ZIP化
zip deployment.zip bootstrap

# Lambda関数の作成
aws lambda create-function \
  --function-name chatbot-backend \
  --runtime provided.al2023 \
  --role $(terraform output -raw backend_role_arn) \
  --handler bootstrap \
  --zip-file fileb://deployment.zip \
  --environment Variables="{KNOWLEDGE_BASE_ID=$KB_ID,MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0}" \
  --timeout 60 \
  --memory-size 512

# API Gateway (オプション)
# Lambda関数URLを有効化
aws lambda create-function-url-config \
  --function-name chatbot-backend \
  --auth-type NONE \
  --cors AllowOrigins="*"
```

#### 5.3 ECS へのデプロイ

```bash
# Dockerfileの作成
cat > Dockerfile <<EOF
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
EOF

# イメージのビルドとプッシュ
docker build -t chatbot-backend .
docker tag chatbot-backend:latest $ECR_REPO:latest
docker push $ECR_REPO:latest

# ECSタスク定義とサービスの作成 (Terraformで管理推奨)
```

### ステップ 6: フロントエンドのデプロイ

#### 6.1 環境設定

```bash
cd ../frontend

# 本番用環境変数を編集
cat > src/environments/environment.prod.ts <<EOF
export const environment = {
  production: true,
  apiUrl: 'https://your-api-domain.com'  # Lambda URL または ALB URL
};
EOF
```

#### 6.2 ビルド

```bash
npm install
npm run build
```

#### 6.3 S3 + CloudFront へのデプロイ

```bash
# S3バケットの作成 (静的ウェブサイトホスティング用)
aws s3 mb s3://my-chatbot-frontend
aws s3 website s3://my-chatbot-frontend --index-document index.html

# ビルド成果物をアップロード
aws s3 sync dist/frontend/ s3://my-chatbot-frontend/

# CloudFrontディストリビューションの作成 (Terraformで管理推奨)
```

## 動作確認

1. フロントエンドのURLにアクセス
2. 接続状態が「Connected」になることを確認
3. サンプル質問を入力:
   - "このシステムについて教えてください"
   - "アップロードしたドキュメントの内容を要約してください"
4. 回答と参照元が表示されることを確認

## トラブルシューティング

### インジェストジョブが失敗する

```bash
# ジョブの詳細を確認
aws bedrock-agent get-ingestion-job \
  --knowledge-base-id $KB_ID \
  --data-source-id $DS_ID \
  --ingestion-job-id <job-id>
```

- IAMロールの権限を確認
- S3バケットのファイル形式を確認

### バックエンドエラー

```bash
# CloudWatch Logsを確認
aws logs tail /aws/lambda/chatbot-backend --follow
```

- AWS認証情報を確認
- Knowledge Base IDが正しいか確認
- Bedrockモデルが有効化されているか確認

### CORS エラー

- バックエンドのCORS設定を確認
- APIのURLが正しいか確認

## クリーンアップ

```bash
# フロントエンドバケットを空にして削除
aws s3 rm s3://my-chatbot-frontend --recursive
aws s3 rb s3://my-chatbot-frontend

# Lambda関数の削除
aws lambda delete-function --function-name chatbot-backend

# Terraformリソースの削除
cd infrastructure
terraform destroy
```

注意: Knowledge BaseとOpenSearch Serverlessの削除には数分かかります。

## セキュリティベストプラクティス

1. **IAMロールの最小権限**
   - 必要最小限の権限のみを付与

2. **VPC配置**
   - 本番環境ではバックエンドとOpenSearchをVPC内に配置

3. **認証・認可**
   - Cognito、Auth0などを統合

4. **APIキー**
   - API Gatewayでレート制限を設定

5. **データ暗号化**
   - 転送中: TLS 1.2以上
   - 保管中: KMSカスタムキーの使用

6. **ログとモニタリング**
   - CloudWatch Logsの有効化
   - AWS CloudTrailの有効化
   - セキュリティアラートの設定
