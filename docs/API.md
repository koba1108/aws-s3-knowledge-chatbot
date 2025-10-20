# API リファレンス

## バックエンド API

ベースURL: `http://localhost:8080` (開発環境)

### 認証

現在の実装では認証は不要です。本番環境では適切な認証機構の実装を推奨します。

---

## エンドポイント

### 1. ヘルスチェック

システムの稼働状態を確認します。

**エンドポイント:** `GET /api/health`

**リクエスト:**
```bash
curl http://localhost:8080/api/health
```

**レスポンス:**
```json
{
  "status": "healthy",
  "time": "2024-01-01T12:00:00Z"
}
```

**ステータスコード:**
- `200 OK`: システムは正常に動作中

---

### 2. チャットメッセージ送信

ユーザーのメッセージを送信し、AIからの応答を取得します。

**エンドポイント:** `POST /api/chat`

**リクエストヘッダー:**
```
Content-Type: application/json
```

**リクエストボディ:**
```json
{
  "message": "AWS S3について教えてください",
  "session_id": "session-12345",
  "knowledge_base_id": "ABC123XYZ"
}
```

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `message` | string | ✓ | ユーザーの質問メッセージ |
| `session_id` | string | - | セッションID（省略時は自動生成） |
| `knowledge_base_id` | string | - | Knowledge Base ID（省略時はデフォルト値を使用） |

**リクエスト例:**
```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "AWS S3の主な機能を説明してください"
  }'
```

**レスポンス（成功時）:**
```json
{
  "response": "AWS S3（Amazon Simple Storage Service）の主な機能は以下の通りです：\n\n1. オブジェクトストレージ: ファイルをオブジェクトとして保存\n2. 耐久性: 99.999999999%（イレブンナイン）の耐久性\n3. スケーラビリティ: 容量制限なし\n4. セキュリティ: 暗号化、アクセス制御\n5. バージョニング: ファイルの変更履歴を保持",
  "session_id": "session-1704110400",
  "sources": [
    {
      "content": "Amazon S3は高い耐久性と可用性を提供するオブジェクトストレージサービスです...",
      "location": {
        "uri": "s3://my-chatbot-kb-prod/aws-s3-guide.pdf"
      }
    }
  ]
}
```

**レスポンスフィールド:**

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `response` | string | AIが生成した回答テキスト |
| `session_id` | string | セッションID（会話の継続に使用） |
| `sources` | array | 回答の根拠となったソースドキュメント（オプション） |
| `sources[].content` | string | 参照されたテキスト内容 |
| `sources[].location.uri` | string | ソースドキュメントのS3 URI |

**レスポンス（エラー時）:**
```json
{
  "response": "",
  "session_id": "session-1704110400",
  "error": "Failed to get response: insufficient permissions"
}
```

**ステータスコード:**
- `200 OK`: リクエスト成功
- `400 Bad Request`: 不正なリクエスト（messageが空など）
- `405 Method Not Allowed`: POSTメソッド以外
- `500 Internal Server Error`: サーバーエラー

---

## エラーレスポンス

APIエラーは以下の形式で返されます：

```json
{
  "error": "エラーメッセージ",
  "session_id": "session-id-if-available"
}
```

### 一般的なエラー

| エラーメッセージ | 原因 | 対処法 |
|----------------|------|--------|
| `Message is required` | messageパラメータが空 | 有効なメッセージを送信 |
| `Invalid request body` | JSONフォーマットエラー | リクエストボディを確認 |
| `Failed to get response: ...` | Bedrock API エラー | AWS認証情報、権限を確認 |
| `Method not allowed` | GET/PUT等の非対応メソッド | POSTメソッドを使用 |

---

## セッション管理

### セッションIDについて

- 同じセッションIDを使用することで、会話のコンテキストが保持されます
- セッションIDを指定しない場合、各リクエストで新しいセッションが作成されます
- セッションは一定時間（Bedrockのデフォルト設定による）後に無効化されます

**会話の継続例:**

```bash
# 最初のメッセージ
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "AWS S3とは何ですか？"}'
# レスポンスからsession_idを取得: "session-1704110400"

# 同じセッションで続きの質問
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "それの料金体系を教えてください",
    "session_id": "session-1704110400"
  }'
```

---

## レート制限

現在の実装ではレート制限はありません。本番環境では以下の対策を推奨します：

- API Gatewayのスロットリング設定
- WAFによるリクエスト制限
- アプリケーションレベルのレート制限実装

---

## CORS設定

開発環境では全オリジンからのアクセスを許可しています：

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: *
```

本番環境では特定のオリジンのみを許可するよう設定してください。

---

## 使用例

### JavaScript (Fetch API)

```javascript
async function sendMessage(message, sessionId) {
  const response = await fetch('http://localhost:8080/api/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      message: message,
      session_id: sessionId
    })
  });
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  return await response.json();
}

// 使用例
sendMessage('AWS S3について教えてください')
  .then(data => {
    console.log('AI Response:', data.response);
    console.log('Session ID:', data.session_id);
    console.log('Sources:', data.sources);
  })
  .catch(error => console.error('Error:', error));
```

### Python

```python
import requests
import json

def send_message(message, session_id=None):
    url = 'http://localhost:8080/api/chat'
    payload = {
        'message': message
    }
    if session_id:
        payload['session_id'] = session_id
    
    response = requests.post(url, json=payload)
    response.raise_for_status()
    return response.json()

# 使用例
result = send_message('AWS S3について教えてください')
print(f"Response: {result['response']}")
print(f"Session ID: {result['session_id']}")

# 同じセッションで続きの質問
result2 = send_message(
    'それの主な用途は？', 
    session_id=result['session_id']
)
print(f"Response: {result2['response']}")
```

### cURL

```bash
# 基本的な使用
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, how are you?"}'

# セッションIDを指定
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "続きを教えてください",
    "session_id": "session-1704110400"
  }'

# Knowledge Base IDを指定
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "質問内容",
    "knowledge_base_id": "CUSTOM_KB_ID"
  }'
```

---

## ベストプラクティス

1. **セッション管理**
   - 会話の継続には同じsession_idを使用
   - セッションIDをクライアント側で保存

2. **エラーハンドリング**
   - ネットワークエラーに対する再試行ロジック
   - タイムアウトの適切な設定

3. **タイムアウト**
   - Bedrock APIの応答には時間がかかる場合があります
   - 60秒程度のタイムアウトを推奨

4. **メッセージ長**
   - 過度に長いメッセージは分割を検討
   - 推奨: 2000文字以内

5. **セキュリティ**
   - 本番環境では認証トークンを実装
   - HTTPS通信の使用
   - 入力のバリデーション
