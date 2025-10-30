package handler

import (
	"aws-s3-knowledge-chatbot/backend/internal/usecase"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	channel, err := h.bedrockAgentRuntimeUsecase.InvokeStream(c, r.SessionID, r.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-channel; ok {
			c.SSEvent("message", gin.H{
				"sessionID": c.Query("session_id"),
				"content":   msg,
			})
			return true
		}
		return false
	})
}
