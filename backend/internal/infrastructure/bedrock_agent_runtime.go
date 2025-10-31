package infrastructure

import (
	"aws-s3-knowledge-chatbot/backend/internal/config"
	"aws-s3-knowledge-chatbot/backend/internal/domain/repository"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	agtypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/samber/lo"
)

type BedrockAgentRuntimeRepository interface {
	RetrieveAndGenerateStream(ctx context.Context, sessionID, inputText string) (*bedrockagentruntime.RetrieveAndGenerateStreamOutput, error)
}

type bedrockAgentRuntimeRepository struct {
	config *config.Config
	client *bedrockagentruntime.Client
}

func NewBedrockAgentRuntimeRepository(
	config *config.Config,
	client *bedrockagentruntime.Client,
) repository.BedrockAgentRuntimeRepository {
	return &bedrockAgentRuntimeRepository{
		config: config,
		client: client,
	}
}

func (r *bedrockAgentRuntimeRepository) RetrieveAndGenerateStream(ctx context.Context, sessionID, inputText string) (*bedrockagentruntime.RetrieveAndGenerateStreamOutput, error) {
	output, err := r.client.RetrieveAndGenerateStream(ctx, &bedrockagentruntime.RetrieveAndGenerateStreamInput{
		SessionId: lo.Ternary(sessionID != "", lo.ToPtr(sessionID), nil),
		Input:     &agtypes.RetrieveAndGenerateInput{Text: &inputText},
		RetrieveAndGenerateConfiguration: &agtypes.RetrieveAndGenerateConfiguration{
			Type: agtypes.RetrieveAndGenerateTypeKnowledgeBase,
			KnowledgeBaseConfiguration: &agtypes.KnowledgeBaseRetrieveAndGenerateConfiguration{
				KnowledgeBaseId: lo.ToPtr(r.config.KnowledgeBaseID),
				ModelArn:        lo.ToPtr(r.config.BedrockModelArn),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call RetrieveAndGenerate: %w", err)
	}
	return output, nil
}

func (r *bedrockAgentRuntimeRepository) RetrieveAndGenerate(ctx context.Context, sessionID, inputText string) (*bedrockagentruntime.RetrieveAndGenerateOutput, error) {
	output, err := r.client.RetrieveAndGenerate(ctx, &bedrockagentruntime.RetrieveAndGenerateInput{
		SessionId: lo.Ternary(sessionID != "", lo.ToPtr(sessionID), nil),
		Input:     &agtypes.RetrieveAndGenerateInput{Text: &inputText},
		RetrieveAndGenerateConfiguration: &agtypes.RetrieveAndGenerateConfiguration{
			Type: agtypes.RetrieveAndGenerateTypeKnowledgeBase,
			KnowledgeBaseConfiguration: &agtypes.KnowledgeBaseRetrieveAndGenerateConfiguration{
				KnowledgeBaseId: lo.ToPtr(r.config.KnowledgeBaseID),
				ModelArn:        lo.ToPtr(r.config.BedrockModelArn),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call RetrieveAndGenerate (non-stream): %w", err)
	}
	return output, nil
}
