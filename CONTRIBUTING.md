# Contributing to AWS S3 Knowledge Chatbot

このプロジェクトへの貢献に興味を持っていただき、ありがとうございます！

## 貢献方法

### バグレポート

バグを見つけた場合は、GitHubのIssuesで報告してください。以下の情報を含めてください：

- 問題の詳細な説明
- 再現手順
- 期待される動作
- 実際の動作
- 環境情報（OS、Go/Node.jsバージョンなど）
- エラーメッセージやログ

### 機能リクエスト

新しい機能の提案は歓迎します。Issueを作成して以下を説明してください：

- 機能の説明
- ユースケース
- 実装案（あれば）

### プルリクエスト

1. このリポジトリをフォーク
2. 新しいブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 開発ガイドライン

### コードスタイル

#### Go
- `gofmt`でフォーマット
- `golint`でチェック
- 明確な変数名を使用
- エラーハンドリングを適切に実装

#### TypeScript/Angular
- Angular Style Guideに従う
- ESLintルールに準拠
- 型を明示的に指定

### テスト

- 新しい機能には必ずテストを追加
- 既存のテストが壊れないことを確認
- カバレッジを維持または向上

```bash
# バックエンドのテスト
cd backend && go test -v ./...

# フロントエンドのテスト
cd frontend && npm test
```

### コミットメッセージ

明確で説明的なコミットメッセージを書いてください：

```
type: 短い説明（50文字以内）

詳細な説明（必要に応じて）

Fixes #issue_number
```

タイプ:
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの動作に影響しない変更（フォーマットなど）
- `refactor`: バグ修正や機能追加ではないコードの変更
- `test`: テストの追加や修正
- `chore`: ビルドプロセスやツールの変更

### ドキュメント

- コードコメントは日本語または英語
- READMEやガイドは日本語で記載
- 新機能には必ずドキュメントを追加

## プロジェクト構造

```
.
├── infrastructure/     # Terraform IaC
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── backend/           # Go API server
│   ├── main.go
│   └── main_test.go
├── frontend/          # Angular SPA
│   └── src/
│       ├── app/
│       └── environments/
└── docs/             # Documentation
    ├── ARCHITECTURE.md
    ├── DEPLOYMENT.md
    └── API.md
```

## 質問

質問がある場合は、GitHubのDiscussionsまたはIssuesで気軽にお尋ねください。

## 行動規範

- 敬意を持って接する
- 建設的なフィードバックを提供
- 多様性を尊重
- コミュニティの成長に貢献

## ライセンス

このプロジェクトに貢献することで、あなたの貢献がMITライセンスの下でライセンスされることに同意したものとみなされます。
