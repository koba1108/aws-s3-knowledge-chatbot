package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
)

type BedrockAgentRuntimeRepository interface {
	RetrieveAndGenerateStream(ctx context.Context, sessionID, inputText string) (*bedrockagentruntime.RetrieveAndGenerateStreamOutput, error)
}
