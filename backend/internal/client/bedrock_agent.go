package client

import (
	"aws-s3-knowledge-chatbot/backend/internal/config"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent/types"
	"github.com/samber/lo"
)

type BedrockAgentClient interface {
	InProgressJobCount(ctx context.Context, limit int32) (int, error)
	StartIngestionJob(ctx context.Context) error
}

type bedrockAgentClient struct {
	client *bedrockagent.Client
	config *config.Config
}

func NewBedrockAgentClient(ctx context.Context, conf *config.Config) (BedrockAgentClient, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.AwsRegion))
	if err != nil {
		return nil, err
	}
	return &bedrockAgentClient{
		client: bedrockagent.NewFromConfig(cfg),
	}, nil
}

func (b *bedrockAgentClient) InProgressJobCount(ctx context.Context, limit int32) (int, error) {
	res, err := b.client.ListIngestionJobs(ctx, &bedrockagent.ListIngestionJobsInput{
		KnowledgeBaseId: aws.String(b.config.KnowledgeBaseID),
		DataSourceId:    aws.String(b.config.DataSourceID),
		MaxResults:      aws.Int32(limit),
	})
	if err != nil {
		return 0, err
	}
	jobs := lo.Filter(res.IngestionJobSummaries, func(item types.IngestionJobSummary, index int) bool {
		return item.Status == types.IngestionJobStatusInProgress || item.Status == types.IngestionJobStatusStarting
	})
	return len(jobs), nil
}

func (b *bedrockAgentClient) StartIngestionJob(ctx context.Context) error {
	res, err := b.client.StartIngestionJob(ctx, &bedrockagent.StartIngestionJobInput{
		KnowledgeBaseId: aws.String(b.config.KnowledgeBaseID),
		DataSourceId:    aws.String(b.config.DataSourceID),
	})
	if err != nil {
		return err
	}
	log.Printf("StartIngestionJob started: %s (status=%s)", lo.FromPtr(res.IngestionJob.IngestionJobId), res.IngestionJob.Status)
	return nil
}
