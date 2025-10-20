# Backend API Server

This is the backend API server written in Go that interfaces with AWS Bedrock Agent Runtime and Knowledge Base.

## Features

- RESTful API for chat interactions
- Integration with AWS Bedrock Knowledge Base
- Session management
- CORS support for frontend
- Health check endpoint

## Prerequisites

- Go 1.20 or later
- AWS credentials configured
- AWS Bedrock Knowledge Base deployed

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Set environment variables:
```bash
export KNOWLEDGE_BASE_ID="your-kb-id"
export MODEL_ID="anthropic.claude-3-sonnet-20240229-v1:0"
export PORT="8080"
export AWS_REGION="us-east-1"
```

## Running

```bash
go run main.go
```

Or build and run:
```bash
go build -o server
./server
```

## API Endpoints

### POST /api/chat

Send a chat message and receive a response from the knowledge base.

Request:
```json
{
  "message": "What is AWS S3?",
  "session_id": "optional-session-id",
  "knowledge_base_id": "optional-kb-id"
}
```

Response:
```json
{
  "response": "AWS S3 is...",
  "session_id": "session-12345",
  "sources": [
    {
      "content": "Source text...",
      "location": {
        "uri": "s3://bucket/file.pdf"
      }
    }
  ]
}
```

### GET /api/health

Health check endpoint.

Response:
```json
{
  "status": "healthy",
  "time": "2024-01-01T00:00:00Z"
}
```

## Testing

```bash
go test -v
```

## Environment Variables

- `KNOWLEDGE_BASE_ID`: AWS Bedrock Knowledge Base ID (required)
- `MODEL_ID`: Bedrock model ID (default: anthropic.claude-3-sonnet-20240229-v1:0)
- `PORT`: Server port (default: 8080)
- `AWS_REGION`: AWS region (default: from AWS config)
