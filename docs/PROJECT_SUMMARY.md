# プロジェクト概要

## 実装内容

このリポジトリには、AWS S3をナレッジベースとして利用するチャットボットアプリケーションの完全な実装が含まれています。

### 技術スタック

- **Infrastructure**: Terraform (AWS)
- **Frontend**: Angular 19
- **Backend**: Go 1.21+
- **AWS Services**: 
  - Amazon S3
  - Amazon Bedrock Knowledge Base
  - Amazon Bedrock Agent Runtime
  - OpenSearch Serverless
  - IAM

## プロジェクト構成

```
aws-s3-knowledge-chatbot/
├── infrastructure/          # Terraformによるインフラ定義
│   ├── main.tf             # メインのリソース定義
│   ├── variables.tf        # 変数定義
│   ├── outputs.tf          # 出力値定義
│   └── README.md           # インフラドキュメント
│
├── backend/                # Goバックエンドサーバー
│   ├── main.go             # メインアプリケーション
│   ├── main_test.go        # ユニットテスト
│   ├── Dockerfile          # Docker イメージ定義
│   └── README.md           # バックエンドドキュメント
│
├── frontend/               # Angularフロントエンド
│   ├── src/
│   │   ├── app/
│   │   │   ├── components/chat/  # チャットUI
│   │   │   ├── services/         # APIサービス
│   │   │   └── models/           # データモデル
│   │   └── environments/         # 環境設定
│   ├── Dockerfile          # Docker イメージ定義
│   ├── nginx.conf          # Nginx設定
│   └── README.md           # フロントエンドドキュメント
│
├── docs/                   # ドキュメント
│   ├── ARCHITECTURE.md     # アーキテクチャ設計
│   ├── DEPLOYMENT.md       # デプロイメントガイド
│   ├── API.md              # API リファレンス
│   └── QUICKSTART.md       # クイックスタート
│
├── .github/
│   └── workflows/
│       └── ci.yml          # CI/CD パイプライン
│
├── docker-compose.yml      # ローカル開発環境
├── Makefile               # 便利なタスクランナー
├── CONTRIBUTING.md        # 貢献ガイド
├── LICENSE                # MITライセンス
└── README.md              # メインドキュメント
```

## 主な機能

### インフラストラクチャ (Terraform)

✅ **S3バケット**
- ドキュメント保管用
- バージョニング有効
- パブリックアクセスブロック

✅ **Bedrock Knowledge Base**
- Titan Embed Text v1による埋め込み
- 固定サイズチャンキング (300トークン、20%オーバーラップ)
- S3データソース統合

✅ **OpenSearch Serverless**
- ベクトルデータベース
- 自動スケーリング
- セキュアな設定

✅ **IAM ロールと権限**
- Knowledge Base用ロール
- バックエンド用ロール
- 最小権限の原則

### バックエンド (Go)

✅ **RESTful API**
- `POST /api/chat` - チャットメッセージ処理
- `GET /api/health` - ヘルスチェック

✅ **Bedrock統合**
- Agent Runtime クライアント
- RetrieveAndGenerate API
- セッション管理

✅ **セキュリティ**
- CORS設定
- エラーハンドリング
- AWS SDK v2使用

✅ **テスト**
- ユニットテスト
- Race detector対応

### フロントエンド (Angular)

✅ **チャットUI**
- リアルタイムメッセージング
- タイピングインジケーター
- ソース表示

✅ **レスポンシブデザイン**
- モダンなUI/UX
- グラデーション背景
- アニメーション

✅ **機能**
- セッション管理
- 接続状態表示
- チャットクリア

✅ **ビルド設定**
- 本番/開発環境分離
- Nginx設定
- Docker対応

## 開発ツール

### Makefile
便利なコマンド集:
```bash
make help              # ヘルプ表示
make install           # 全依存関係インストール
make build             # 全ビルド
make test              # 全テスト実行
make terraform-apply   # インフラデプロイ
```

### Docker Compose
ローカル開発環境:
```bash
docker-compose up      # 全サービス起動
```

### CI/CD (GitHub Actions)
- バックエンドテスト
- フロントエンドテスト
- Terraform検証
- セキュリティスキャン
- Dockerイメージビルド

## ドキュメント

### ユーザー向け
- **README.md**: プロジェクト概要とセットアップ
- **QUICKSTART.md**: 5分で動かすガイド
- **DEPLOYMENT.md**: 本番環境デプロイ

### 開発者向け
- **ARCHITECTURE.md**: システム設計詳細
- **API.md**: API仕様書
- **CONTRIBUTING.md**: 貢献ガイドライン

## セキュリティ

✅ **脆弱性スキャン**
- CodeQL: 脆弱性なし
- Go依存関係: 脆弱性なし
- npm監査: クリーン

✅ **ベストプラクティス**
- IAM最小権限
- データ暗号化
- セキュアな通信
- 入力検証

## テスト結果

### バックエンド
```
=== RUN   TestHandleHealth
--- PASS: TestHandleHealth (0.00s)
=== RUN   TestHandleChatInvalidMethod
--- PASS: TestHandleChatInvalidMethod (0.00s)
=== RUN   TestHandleChatEmptyMessage
--- PASS: TestHandleChatEmptyMessage (0.00s)
PASS
```

### フロントエンド
```
Application bundle generation complete.
Initial chunk files   | Names     | Raw size
main-ZTLAD3P5.js     | main      | 294.63 kB
polyfills-5CFQRCPP.js| polyfills |  34.59 kB
styles-ZLPCNYXJ.css  | styles    | 161 bytes
```

## 次のステップ

1. **クイックスタート**: `docs/QUICKSTART.md` を参照
2. **ローカルテスト**: Docker Composeで起動
3. **本番デプロイ**: `docs/DEPLOYMENT.md` を参照
4. **カスタマイズ**: アーキテクチャドキュメントを確認

## 技術的な特徴

### RAG (Retrieval-Augmented Generation)
- ベクトル検索による関連情報取得
- LLMによる自然な回答生成
- ソース引用機能

### スケーラビリティ
- ステートレス設計
- 水平スケーリング対応
- サーバーレス互換

### 保守性
- 明確なディレクトリ構造
- 包括的なドキュメント
- 自動テストとCI/CD

## ライセンス

MIT License - 詳細は `LICENSE` ファイルを参照

## サポート

- GitHub Issues: バグ報告・機能リクエスト
- GitHub Discussions: 質問・議論
- Pull Requests: 貢献歓迎

---

このプロジェクトは、AWS Bedrockを活用した最新のRAGアプリケーションの
リファレンス実装として設計されています。
