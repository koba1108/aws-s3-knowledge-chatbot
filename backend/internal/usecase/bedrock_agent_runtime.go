package usecase

import (
	"aws-s3-knowledge-chatbot/backend/internal/domain/repository"
	"context"
	"fmt"
	"log"

	atypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/samber/lo"
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
	log.Printf("RetrieveAndGenerateStream SessionId: %+v\n", res.SessionId)
	log.Printf("RetrieveAndGenerateStream ResultMetadata: %+v\n", res.ResultMetadata)

	stream := res.GetStream()
	if stream == nil {
		return nil, fmt.Errorf("nil stream returned")
	}
	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream error: %w", err)
	}

	outputChan := make(chan string)

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
				outputChan <- lo.FromPtr(e.Value.Text)
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberCitation:
				log.Printf("[stream] citation: %+v\n", e.Value)
			case *atypes.RetrieveAndGenerateStreamResponseOutputMemberGuardrail:
				log.Printf("[stream] guardrail: %+v\n", e.Value)
			default:
				log.Printf("[stream] unknown event: %T %+v\n", e, e)
			}
		}
		log.Printf("Finished processing RetrieveAndGenerateStream events (count=%d)", cnt)
	}()

	return outputChan, nil
}
