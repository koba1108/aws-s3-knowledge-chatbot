# 開発用（ホットリロード）イメージ
FROM --platform=linux/arm64 golang:1.25.0-alpine

# HTTPS通信用（AWS SDK用）＋ビルドに必要な最低限
RUN apk add --no-cache ca-certificates git build-base

WORKDIR /app

# 依存キャッシュを効かせる
COPY go.mod go.sum ./
RUN go mod download

# Air を導入（ホットリロード）
RUN go install github.com/air-verse/air@latest

# Air の設定ファイルはボリュームで渡す想定（無い場合は Air がデフォルト挙動）
CMD ["air", "-c", "/app/backend/cmd/bedrock_agent_runtime/.air.toml"]
