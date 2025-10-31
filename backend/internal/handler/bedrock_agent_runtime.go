package handler

import (
	"aws-s3-knowledge-chatbot/backend/internal/transport/http/sse"
	"aws-s3-knowledge-chatbot/backend/internal/usecase"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BedrockAgentRuntimeHandler interface {
	Ping(ctx *gin.Context)
	InvokeStream(ctx *gin.Context)
}

type bedrockAgentRuntimeHandler struct {
	bedrockAgentRuntimeUsecase usecase.BedrockAgentRuntimeUsecase
}

func NewBedrockAgentRuntimeHandler(bedrockAgentRuntimeUsecase usecase.BedrockAgentRuntimeUsecase) BedrockAgentRuntimeHandler {
	return &bedrockAgentRuntimeHandler{
		bedrockAgentRuntimeUsecase: bedrockAgentRuntimeUsecase,
	}
}

func (h *bedrockAgentRuntimeHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339Nano),
	})
}

func (h *bedrockAgentRuntimeHandler) InvokeStream(c *gin.Context) {
	type req struct {
		SessionID string `json:"session_id"`
		Query     string `json:"query" binding:"required"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	em := sse.NewEmitter(c)
	stopHeartbeat := em.StartHeartbeat(10 * time.Second)
	defer stopHeartbeat()

	reqCtx := c.Request.Context()
	srvCtx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	// 片方が Done になったら終わるように束ねる
	ctx, stop := context.WithCancelCause(srvCtx)
	defer stop(nil)
	go func() {
		select {
		case <-reqCtx.Done():
			stop(reqCtx.Err()) // クライアント切断を優先的に伝播
		case <-srvCtx.Done():
			// タイムアウト/サーバ事情
		}
	}()

	ch, err := h.bedrockAgentRuntimeUsecase.InvokeStream(ctx, r.SessionID, r.Query)
	if err != nil {
		_ = em.EmitError(err.Error(), sse.WithSessionID(r.SessionID))
		return
	}

	sentStarted := false
	c.Stream(func(_ io.Writer) bool {
		select {
		case <-ctx.Done():
			_ = em.EmitMessageEnd(sse.FinishError, sse.WithSessionID(r.SessionID))
			return false
		case msg, ok := <-ch:
			if !ok {
				_ = em.EmitMessageEnd(sse.FinishCompleted, sse.WithSessionID(r.SessionID))
				return false
			}
			if !sentStarted {
				_ = em.EmitMessageStart(sse.RoleAssistant, sse.WithSessionID(r.SessionID))
				sentStarted = true
			}
			_ = em.EmitMessageDelta(msg, sse.WithSessionID(r.SessionID))
			return true
		}
	})
}
