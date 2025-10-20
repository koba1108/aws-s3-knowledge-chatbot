# クイックスタートガイド

このガイドでは、ローカル環境でAWS S3 Knowledge Chatbotを素早くセットアップして動作確認する手順を説明します。

## 前提条件

- AWS アカウント
- AWS CLI がインストール済みで、認証情報が設定済み
- Go 1.20 以上
- Node.js 18 以上
- Terraform 1.0 以上

## 5分でセットアップ

### 1. リポジトリのクローン

```bash
git clone https://github.com/koba1108/aws-s3-knowledge-chatbot.git
cd aws-s3-knowledge-chatbot
```

### 2. インフラのデプロイ (約3-5分)

```bash
cd infrastructure
terraform init
terraform apply -auto-approve
```

出力されたKnowledge Base IDをメモしておきます：

```bash
export KB_ID=$(terraform output -raw knowledge_base_id)
export DS_ID=$(terraform output -raw data_source_id)
export BUCKET_NAME=$(terraform output -raw s3_bucket_name)
```

### 3. サンプルドキュメントのアップロード

```bash
# サンプルFAQファイルを作成
cat > /tmp/sample-faq.txt << 'EOF'
# AWS S3 Knowledge Chatbot FAQ

## Q1: このシステムは何ができますか？
A1: このシステムは、S3バケットにアップロードされたドキュメントを知識ベースとして、
ユーザーの質問に対してAIが回答を生成します。Amazon Bedrockを活用したRAGシステムです。

## Q2: どのようなドキュメント形式に対応していますか？
A2: PDF、TXT、DOCX、HTML、Markdownなど、様々なテキストベースのドキュメントに対応しています。

## Q3: セキュリティはどうなっていますか？
A3: S3バケットはパブリックアクセスがブロックされ、IAMロールによる適切な権限管理が
実装されています。データは転送中・保管中ともに暗号化されます。

## Q4: コストはどのくらいかかりますか？
A4: 主なコストは以下です：
- Amazon Bedrock: APIコールとトークン数に応じた従量課金
- OpenSearch Serverless: OCU（OpenSearch Compute Unit）課金
- S3: ストレージ使用量に応じた課金
小規模な利用であれば月額10-30ドル程度です。

## Q5: カスタマイズは可能ですか？
A5: はい、使用するLLMモデル、埋め込みモデル、チャンクサイズなど、
様々なパラメータをカスタマイズ可能です。
EOF

# S3にアップロード
aws s3 cp /tmp/sample-faq.txt s3://$BUCKET_NAME/sample-faq.txt
```

### 4. データの同期 (約2-3分)

```bash
# インジェストジョブを開始
aws bedrock-agent start-ingestion-job \
  --knowledge-base-id $KB_ID \
  --data-source-id $DS_ID \
  --region $(terraform output -raw aws_region || echo "us-east-1")

# 完了を待つ（約2-3分）
echo "インジェスト中...完了するまでお待ちください"
sleep 180
```

### 5. バックエンドの起動

別のターミナルで：

```bash
cd ../backend
export KNOWLEDGE_BASE_ID=$KB_ID
export AWS_REGION=$(cd ../infrastructure && terraform output -raw aws_region || echo "us-east-1")
go run main.go
```

出力例：
```
2024/01/01 12:00:00 Server starting on port 8080
2024/01/01 12:00:00 Knowledge Base ID: ABC123XYZ
```

### 6. フロントエンドの起動

さらに別のターミナルで：

```bash
cd frontend
npm install  # 初回のみ
npm start
```

ブラウザが自動的に開きます（または http://localhost:4200 にアクセス）。

## 動作確認

### テスト質問

フロントエンドのチャット画面で以下の質問を試してください：

1. **基本的な質問**
   ```
   このシステムは何ができますか？
   ```

2. **詳細な質問**
   ```
   どのようなドキュメント形式に対応していますか？
   ```

3. **技術的な質問**
   ```
   セキュリティについて教えてください
   ```

4. **コストに関する質問**
   ```
   このシステムのコストはどのくらいですか？
   ```

### 期待される動作

- AIが質問に対して適切な回答を返す
- 回答の下に「Sources」として参照元のドキュメント（sample-faq.txt）が表示される
- セッションが保持され、会話のコンテキストが維持される

## トラブルシューティング

### バックエンドが起動しない

```bash
# AWS認証情報を確認
aws sts get-caller-identity

# Knowledge Base IDを確認
echo $KNOWLEDGE_BASE_ID
```

### フロントエンドから接続できない

1. バックエンドが起動しているか確認
2. ブラウザのコンソールでエラーを確認
3. CORS設定を確認

### AIが応答しない

```bash
# インジェストジョブの状態を確認
aws bedrock-agent list-ingestion-jobs \
  --knowledge-base-id $KB_ID \
  --data-source-id $DS_ID \
  --max-results 5
```

状態が「COMPLETE」であることを確認してください。

### Bedrockモデルへのアクセスエラー

AWS Bedrockのモデルアクセスを有効化：

1. AWSマネジメントコンソールを開く
2. Bedrockサービスに移動
3. 左メニューから「Model access」を選択
4. 「Manage model access」をクリック
5. 以下のモデルを有効化：
   - Anthropic Claude 3 Sonnet
   - Amazon Titan Embeddings G1 - Text

## 追加のドキュメントをアップロード

```bash
# 独自のドキュメントをアップロード
aws s3 cp your-document.pdf s3://$BUCKET_NAME/

# 再度インジェストを実行
aws bedrock-agent start-ingestion-job \
  --knowledge-base-id $KB_ID \
  --data-source-id $DS_ID
```

## クリーンアップ

テスト終了後、リソースを削除：

```bash
cd infrastructure
terraform destroy -auto-approve
```

## 次のステップ

- [デプロイメントガイド](./DEPLOYMENT.md) - 本番環境への展開
- [アーキテクチャドキュメント](./ARCHITECTURE.md) - システム設計の詳細
- [API リファレンス](./API.md) - API仕様の詳細

## サポート

問題が発生した場合は、GitHubのIssuesでお知らせください。
