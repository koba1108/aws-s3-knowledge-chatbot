package client

import (
	"aws-s3-knowledge-chatbot/backend/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/smithy-go/logging"
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
	ac.ClientLogMode = aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogRequestWithBody | aws.LogResponseWithBody
	return bedrockagentruntime.NewFromConfig(ac, func(o *bedrockagentruntime.Options) {
		o.Logger = logging.NewStandardLogger(log.Writer())
	}), nil
}

func NewBedrockAgentRuntimeClientMust(config *config.Config) *bedrockagentruntime.Client {
	client, err := NewBedrockAgentRuntimeClient(config)
	if err != nil {
		panic(err)
	}
	return client
}
