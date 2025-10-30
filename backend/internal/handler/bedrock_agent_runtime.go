package handler

import (
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

	// --- SSE headers ---
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// Nginxなどでバッファされないように
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// --- Context: クライアント切断 + サーバ側タイムアウトの両立 ---
	reqCtx := c.Request.Context()
	srvCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	// --- 呼び出し ---
	ch, err := h.bedrockAgentRuntimeUsecase.InvokeStream(ctx, r.SessionID, r.Query)
	if err != nil {
		// SSEとしてエラーを返す（JSON 500ではSSE接続が閉じるだけで原因が見えない）
		c.SSEvent("error", gin.H{"message": err.Error()})
		c.Writer.Flush()
		return
	}

	// --- ストリーミングループ ---
	heartbeat := time.NewTicker(10 * time.Second)
	defer heartbeat.Stop()

	c.Stream(func(io.Writer) bool {
		select {
		case <-ctx.Done():
			// 終了イベント（任意）
			c.SSEvent("done", gin.H{
				"sessionID": r.SessionID,
				"reason":    ctx.Err().Error(),
			})
			return false

		case msg, ok := <-ch:
			if !ok {
				// チャネルが閉じられた＝サーバ側完了
				c.SSEvent("done", gin.H{
					"sessionID": r.SessionID,
					"reason":    "completed",
				})
				return false
			}
			c.SSEvent("message", gin.H{
				"sessionID": r.SessionID,
				"content":   msg,
			})
			return true

		case <-heartbeat.C:
			// 無音対策でハートビート
			// SSEのコメント行はクライアントからは無視されるが接続維持に効く
			_, _ = c.Writer.Write([]byte(":heartbeat\n\n"))
			c.Writer.Flush()
			return true
		}
	})
}
