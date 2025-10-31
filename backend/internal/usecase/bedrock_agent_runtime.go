package usecase

import (
	"aws-s3-knowledge-chatbot/backend/internal/domain/repository"
	"aws-s3-knowledge-chatbot/backend/internal/transport/http/sse"
	"context"
	"fmt"
	"log"

	atypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/samber/lo"
)

type BedrockAgentRuntimeUsecase interface {
	InvokeStream(ctx context.Context, sessionId, query string) (<-chan sse.AIEvent, error)
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

func (u *bedrockAgentRuntimeUsecase) InvokeStream(ctx context.Context, sessionId, query string) (<-chan sse.AIEvent, error) {
	res, err := u.bedrockAgentRuntimeRepository.RetrieveAndGenerateStream(ctx, sessionId, query)
	if err != nil {
		return nil, err
	}

	stream := res.GetStream()
	if stream == nil {
		return nil, fmt.Errorf("nil stream returned")
	}
	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream error: %w", err)
	}

	outputChan := make(chan sse.AIEvent)

	go func() {
		defer func() {
			// 明示クローズ＆出力チャネルを閉じる
			_ = stream.Close()
			close(outputChan)
		}()

		cnt := 0
		for ev := range stream.Events() {
			cnt++
			switch e := ev.(type) {
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberOutput:
				if e.Value.Text != nil {
					outputChan <- sse.NewAssistantDelta(lo.FromPtr(e.Value.Text))
				}
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberCitation:
				outputChan <- sse.NewAIMessageCitation(lo.Map(e.Value.RetrievedReferences, func(ref atypes.RetrievedReference, _ int) sse.CitationReference {
					return sse.CitationReference{
						Text:   lo.FromPtr(ref.Content.Text),
						Source: lo.FromPtr(ref.Location.S3Location.Uri),
					}
				}))
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberGuardrail:
				log.Printf("[stream] guardrail: %+v\n", e.Value)
			default:
				log.Printf("[stream] unknown event: %T %+v\n", e, e)
			}
		}
	}()

	return outputChan, nil
}
