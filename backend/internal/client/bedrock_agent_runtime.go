package client

import (
	"aws-s3-knowledge-chatbot/backend/internal/config"
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
)

type BedrockAgentRuntime struct {
	config *config.Config
	Client *bedrockagentruntime.Client
}

func NewBedrockAgentRuntimeClient(config *config.Config) (*bedrockagentruntime.Client, error) {
	ac, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(config.AwsRegion))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	return bedrockagentruntime.NewFromConfig(ac), nil
}

func NewBedrockAgentRuntimeClientMust(config *config.Config) *bedrockagentruntime.Client {
	client, err := NewBedrockAgentRuntimeClient(config)
	if err != nil {
		panic(err)
	}
	return client
}
