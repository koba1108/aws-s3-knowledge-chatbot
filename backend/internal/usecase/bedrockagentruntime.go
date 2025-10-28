package usecase

import (
	"aws-s3-knowledge-chatbot/backend/internal/domain/repository"
	"context"

	atypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

type BedrockAgentRuntimeUsecase interface {
	InvokeStream(ctx context.Context, sessionId, query string) (<-chan string, error)
}

type bedrockAgentRuntimeUsecase struct {
	bedrockAgentRuntimeRepository repository.BedrockAgentRuntimeRepository
}

func NewBedrockAgentRuntimeUsecase(
	bedrockAgentRuntimeRepository repository.BedrockAgentRuntimeRepository,
) BedrockAgentRuntimeUsecase {
	return &bedrockAgentRuntimeUsecase{
		bedrockAgentRuntimeRepository: bedrockAgentRuntimeRepository,
	}
}

func (u *bedrockAgentRuntimeUsecase) InvokeStream(ctx context.Context, sessionId, query string) (<-chan string, error) {
	res, err := u.bedrockAgentRuntimeRepository.RetrieveAndGenerateStream(ctx, sessionId, query)
	if err != nil {
		return nil, err
	}

	outputChan := make(chan string)
	go func() {
		defer close(outputChan)
		for event := range res.GetStream().Events() {
			switch e := event.(type) {
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberCitation:
				// todo: 引用情報の処理が必要な場合はここに実装
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberGuardrail:
				// todo: ガードレール情報の処理が必要な場合はここに実装
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberOutput:
				if e.Value.Text != nil {
					outputChan <- *e.Value.Text
				}
			}
		}
	}()

	return outputChan, nil
}
