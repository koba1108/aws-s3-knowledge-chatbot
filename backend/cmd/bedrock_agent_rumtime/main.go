package main

import (
	"aws-s3-knowledge-chatbot/backend/internal/client"
	"aws-s3-knowledge-chatbot/backend/internal/config"
	"aws-s3-knowledge-chatbot/backend/internal/handler"
	"aws-s3-knowledge-chatbot/backend/internal/infrastructure"
	"aws-s3-knowledge-chatbot/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfigMust()
	bedrockAgentRuntimeClient := client.NewBedrockAgentRuntimeClientMust(cfg)
	bedrockAgentRuntimeRepository := infrastructure.NewBedrockAgentRuntimeRepository(cfg, bedrockAgentRuntimeClient)
	bedrockAgentRuntimeUsecase := usecase.NewBedrockAgentRuntimeUsecase(bedrockAgentRuntimeRepository)
	bh := handler.NewBedrockAgentRuntimeHandler(bedrockAgentRuntimeUsecase)

	e := gin.Default()
	_ = e.SetTrustedProxies(nil)

	// bedrockAgentRuntimeで必須なエンドポイントを設定
	e.GET("/ping", bh.Ping)
	e.GET("/invocations", bh.InvokeStream)

	if err := e.Run(cfg.GetAddress()); err != nil {
		panic(err)
	}
}
