# ====== Build stage ======
FROM --platform=linux/arm64 golang:1.25.0-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY .. .

# AgentCore Runtime 対応: linux/arm64, static build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -trimpath -ldflags="-s -w" \
    -o /app/server ./backend/cmd/bedrock_agent_rumtime

# ====== Runtime stage ======
FROM --platform=linux/arm64 alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
