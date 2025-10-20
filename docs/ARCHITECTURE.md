# アーキテクチャドキュメント

## システム概要

AWS S3 Knowledge Chatbotは、Amazon Bedrockを活用したRAG (Retrieval-Augmented Generation) システムです。S3に保存されたドキュメントを知識ベースとして、ユーザーの質問に対して正確な回答を提供します。

## コンポーネント構成

### 1. データレイヤー

#### Amazon S3
- **役割**: ソースドキュメントの保管
- **対応フォーマット**: PDF, TXT, DOCX, HTML, MD など
- **バージョニング**: 有効化
- **セキュリティ**: パブリックアクセスブロック有効

#### OpenSearch Serverless
- **役割**: ベクトルデータベース
- **インデックス**: bedrock-knowledge-base-index
- **フィールド構成**:
  - ベクトルフィールド: embeddings
  - テキストフィールド: AMAZON_BEDROCK_TEXT_CHUNK
  - メタデータフィールド: AMAZON_BEDROCK_METADATA

### 2. AI/MLレイヤー

#### Amazon Bedrock Knowledge Base
- **埋め込みモデル**: Amazon Titan Embed Text v1
- **チャンク戦略**: 固定サイズ (300トークン、20%オーバーラップ)
- **データソース**: S3バケット

#### Amazon Bedrock Agent Runtime
- **LLMモデル**: Anthropic Claude 3 Sonnet
- **API**: RetrieveAndGenerate
- **セッション管理**: 会話コンテキストの保持

### 3. アプリケーションレイヤー

#### バックエンド (Go)
- **フレームワーク**: net/http (標準ライブラリ)
- **AWS SDK**: aws-sdk-go-v2
- **エンドポイント**:
  - `POST /api/chat`: チャットメッセージの処理
  - `GET /api/health`: ヘルスチェック
- **CORS**: 全オリジン許可 (開発用)

#### フロントエンド (Angular 19)
- **コンポーネント**: 
  - ChatComponent: メインチャットUI
- **サービス**:
  - ChatService: API通信
- **機能**:
  - リアルタイムチャット
  - セッション管理
  - ソース表示
  - 接続状態モニタリング

## データフロー

```
[ユーザー] 
    ↓ 質問入力
[Angular Frontend]
    ↓ HTTP POST /api/chat
[Go Backend]
    ↓ RetrieveAndGenerate API
[Bedrock Agent Runtime]
    ↓ ベクトル検索
[Knowledge Base] → [OpenSearch Serverless]
    ↓ 関連ドキュメント取得
[S3 Bucket]
    ↓ コンテキスト付きプロンプト
[Claude 3 Sonnet]
    ↓ 生成された回答
[Go Backend]
    ↓ JSON レスポンス
[Angular Frontend]
    ↓ 画面表示
[ユーザー]
```

## セキュリティ

### IAMロール

1. **bedrock-kb-role**
   - S3バケットへの読み取り専用アクセス
   - OpenSearch Serverlessへの書き込みアクセス

2. **backend-role**
   - Bedrock Agent Runtimeの呼び出し権限
   - CloudWatch Logsへの書き込み権限

### ネットワーク

- OpenSearch Serverless: パブリックアクセス (開発用)
- 本番環境: VPC内配置を推奨

### データ暗号化

- S3: サーバーサイド暗号化 (SSE-S3)
- OpenSearch: AWS管理キーで暗号化
- 通信: TLS 1.2以上

## スケーラビリティ

### フロントエンド
- 静的ファイル配信 (CDN推奨)
- サーバーレス対応

### バックエンド
- ステートレス設計
- 水平スケーリング可能
- Lambda/ECS/EKSに対応

### データストア
- OpenSearch Serverless: 自動スケーリング
- S3: 無制限ストレージ

## モニタリング

### メトリクス
- Bedrock API呼び出し数
- レスポンス時間
- エラー率
- セッション数

### ログ
- CloudWatch Logs: バックエンドログ
- X-Ray: 分散トレーシング (オプション)

## コスト最適化

1. **Bedrock**: オンデマンド課金
   - 入力トークン数削減
   - 適切なチャンクサイズ設定

2. **OpenSearch Serverless**: OCU課金
   - 最小OCU設定の調整

3. **S3**: ストレージ使用量
   - ライフサイクルポリシーの設定
   - Intelligent-Tieringの活用

## 今後の拡張

1. **マルチモーダル対応**
   - 画像、音声ファイルのサポート

2. **ファインチューニング**
   - カスタムモデルの利用

3. **エンタープライズ機能**
   - ユーザー認証・認可
   - テナント分離
   - 監査ログ

4. **高度な検索**
   - メタデータフィルタリング
   - ハイブリッド検索
