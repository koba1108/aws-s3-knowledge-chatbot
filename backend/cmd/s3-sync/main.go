package main

import (
	"aws-s3-knowledge-chatbot/backend/internal/client"
	"aws-s3-knowledge-chatbot/backend/internal/config"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) error {
	// 設定読み込み
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println("Failed to load config:", err)
		return err
	}
	bedrockAgent, err := client.NewBedrockAgentClient(ctx, cfg)
	if err != nil {
		log.Println("Failed to create Bedrock Agent client:", err)
		return err
	}
	// 重複チェック
	inProgressCount, err := bedrockAgent.InProgressJobCount(ctx, 1)
	if err != nil {
		log.Println("Failed to get in-progress job count:", err)
		return err
	}
	if inProgressCount > 0 {
		log.Println("Ingestion job is already in progress. Skipping new job start.")
		return nil
	}
	// ジョブ開始
	if err := bedrockAgent.StartIngestionJob(ctx); err != nil {
		log.Println("Failed to start ingestion job:", err)
		return err
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
