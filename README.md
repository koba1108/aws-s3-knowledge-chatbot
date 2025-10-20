# AWS S3 Knowledge Chatbot

このリポジトリは、AWS S3をナレッジベースとして利用するチャットボットアプリケーションです。

## アーキテクチャ

- **Infrastructure**: Terraform
- **Frontend**: Angular 19
- **Backend**: Go (Golang)
- **AWS Services**: 
  - Amazon S3 (ナレッジベースのドキュメント保管)
  - Amazon Bedrock Knowledge Base (ベクトル検索)
  - Amazon Bedrock Agent Runtime (チャット応答生成)
  - OpenSearch Serverless (ベクトルストア)

## 仕組み

1. S3バケットにドキュメント（PDF、テキストなど）をアップロード
2. Bedrock Knowledge Baseがドキュメントをインジェストし、OpenSearch Serverlessにベクトル化して保存
3. ユーザーがフロントエンドから質問を入力
4. バックエンドAPIがBedrock Agent Runtimeを呼び出し、Knowledge Baseから関連情報を取得
5. 取得した情報をもとにLLM（Claude 3 Sonnet）が回答を生成
6. フロントエンドに回答と参照元を表示

## プロジェクト構成

```
.
├── infrastructure/     # Terraform設定ファイル
│   ├── main.tf        # メインのインフラ定義
│   ├── variables.tf   # 変数定義
│   └── outputs.tf     # 出力値
├── backend/           # Goバックエンド
│   ├── main.go        # メインアプリケーション
│   └── main_test.go   # テスト
├── frontend/          # Angularフロントエンド
│   └── src/
│       ├── app/
│       │   ├── components/chat/  # チャットコンポーネント
│       │   ├── services/         # API通信サービス
│       │   └── models/           # データモデル
│       └── environments/         # 環境設定
└── docs/              # ドキュメント
```

## セットアップ

### 1. インフラストラクチャのデプロイ

```bash
cd infrastructure
terraform init
terraform plan
terraform apply
```

デプロイ後、以下の出力値を記録してください：
- `knowledge_base_id`: バックエンドで使用
- `s3_bucket_name`: ドキュメントのアップロード先

### 2. S3バケットにドキュメントをアップロード

```bash
aws s3 cp your-document.pdf s3://[bucket-name]/
```

### 3. Knowledge Baseのデータソースを同期

```bash
aws bedrock-agent start-ingestion-job \
  --knowledge-base-id [your-kb-id] \
  --data-source-id [your-data-source-id]
```

### 4. バックエンドの起動

```bash
cd backend
cp .env.example .env
# .envファイルを編集してKNOWLEDGE_BASE_IDを設定
export KNOWLEDGE_BASE_ID=[your-kb-id]
export AWS_REGION=us-east-1
go run main.go
```

バックエンドは `http://localhost:8080` で起動します。

### 5. フロントエンドの起動

```bash
cd frontend
npm install
npm start
```

フロントエンドは `http://localhost:4200` で起動します。

## 使用方法

1. ブラウザで `http://localhost:4200` にアクセス
2. チャット画面が表示されます
3. 質問を入力して送信
4. AIが回答を生成し、参照元のドキュメントも表示されます

## 開発

### バックエンドのテスト

```bash
cd backend
go test -v
```

### フロントエンドのビルド

```bash
cd frontend
npm run build
```

ビルド成果物は `frontend/dist/` に出力されます。

## 本番環境へのデプロイ

### バックエンド

- AWS Lambda、ECS、EC2などにデプロイ可能
- 環境変数でKnowledge Base IDとAWSリージョンを設定

### フロントエンド

- S3 + CloudFront、またはその他の静的ホスティングサービス
- `environment.prod.ts` でバックエンドAPIのURLを設定

## 必要な権限

### Bedrock Knowledge Base用
- S3バケットへの読み取りアクセス
- OpenSearch Serverlessへのアクセス

### バックエンド用
- Bedrock Agent Runtimeの呼び出し権限
- CloudWatch Logsへの書き込み権限

## トラブルシューティング

### Knowledge Baseが応答を返さない
- データソースの同期が完了しているか確認
- S3バケットにドキュメントがアップロードされているか確認

### 接続エラー
- バックエンドが起動しているか確認
- CORS設定が正しいか確認
- AWS認証情報が設定されているか確認

## ライセンス

このプロジェクトはオープンソースです。

## 参考リンク

- [Amazon Bedrock Documentation](https://docs.aws.amazon.com/bedrock/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [Angular Documentation](https://angular.io/docs)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
