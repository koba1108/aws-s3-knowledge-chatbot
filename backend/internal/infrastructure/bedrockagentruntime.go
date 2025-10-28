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

const claudeSonnet45ModelArn = "anthropic.claude-sonnet-4-5-20250929-v1:0"

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
				ModelArn:        lo.ToPtr(claudeSonnet45ModelArn),
				GenerationConfiguration: &agtypes.GenerationConfiguration{
					// 必要に応じて設定を追加
				},
				RetrievalConfiguration: &agtypes.KnowledgeBaseRetrievalConfiguration{
					// 必要に応じて設定を追加
				},
			},
			ExternalSourcesConfiguration: &agtypes.ExternalSourcesRetrieveAndGenerateConfiguration{
				Sources: []agtypes.ExternalSource{{
					S3Location: &agtypes.S3ObjectDoc{
						Uri: lo.ToPtr(fmt.Sprintf("s3://%s/%s", r.config.KnowledgeS3Bucket, r.config.KnowledgeS3Prefix)),
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call RetrieveAndGenerate: %w", err)
	}
	return output, nil
}
