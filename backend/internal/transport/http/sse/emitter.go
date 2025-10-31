package sse

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	headerContentType     = "text/event-stream"
	headerCacheControl    = "no-cache"
	headerConnection      = "keep-alive"
	headerXAccelBuffering = "no" // nginx等のバッファ無効化
)

// SetupSSEHeaders sets standard SSE headers on the response.
func SetupSSEHeaders(c *gin.Context) {
	h := c.Writer.Header()
	h.Set("Content-Type", headerContentType)
	h.Set("Cache-Control", headerCacheControl)
	h.Set("Connection", headerConnection)
	h.Set("X-Accel-Buffering", headerXAccelBuffering)
}

// Emitter encapsulates SSE writing & flushing.
type Emitter struct {
	w       http.ResponseWriter
	bw      *bufio.Writer
	flusher http.Flusher
	ctx     context.Context
}

// NewEmitter creates an SSE emitter and writes initial headers.
func NewEmitter(c *gin.Context) *Emitter {
	SetupSSEHeaders(c)
	w := c.Writer
	fl, _ := w.(http.Flusher)
	return &Emitter{
		w:       w,
		bw:      bufio.NewWriterSize(w, 32*1024),
		flusher: fl,
		ctx:     c.Request.Context(),
	}
}

// Emit writes an SSE event with JSON-encoded data.
func (e *Emitter) Emit(event string, v any, opts ...EventOption) error {
	if base, ok := v.(interface{ GetBase() *AIBaseEvent }); ok {
		b := base.GetBase()
		for _, opt := range opts {
			opt(b)
		}
	}

	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("sse marshal: %w", err)
	}

	if _, err := e.bw.WriteString("event: " + event + "\n"); err != nil {
		return err
	}
	if _, err := e.bw.WriteString("data: "); err != nil {
		return err
	}
	if _, err := e.bw.Write(b); err != nil {
		return err
	}
	if _, err := e.bw.WriteString("\n\n"); err != nil {
		return err
	}
	return e.Flush()
}

// Flush flushes buffered data to the client.
func (e *Emitter) Flush() error {
	if err := e.bw.Flush(); err != nil {
		return err
	}
	if e.flusher != nil {
		e.flusher.Flush()
	}
	return nil
}

// Comment sends an SSE comment line (useful for heartbeats).
func (e *Emitter) Comment(text string) error {
	if _, err := e.bw.WriteString(":" + text + "\n\n"); err != nil {
		return err
	}
	return e.Flush()
}

// StartHeartbeat emits periodic ping comments until ctx is done.
func (e *Emitter) StartHeartbeat(interval time.Duration) (stop func()) {
	if interval <= 0 {
		interval = 20 * time.Second
	}
	t := time.NewTicker(interval)
	done := make(chan struct{})
	go func() {
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-e.ctx.Done():
				return
			case <-t.C:
				_ = e.Comment("ping")
			}
		}
	}()
	return func() { close(done) }
}

// Close finalizes the stream (optional place for trailers).
func (e *Emitter) Close() error {
	return e.Flush()
}

func (a *AIBaseEvent) GetBase() *AIBaseEvent { return a }

// EmitMessageStart sends "message.start".
func (e *Emitter) EmitMessageStart(role Role, opts ...EventOption) error {
	ev := NewAIMessageStart(role, opts...)
	return e.Emit(string(EventMessageStart), ev, opts...)
}

// EmitMessageDelta sends "message.delta".
func (e *Emitter) EmitMessageDelta(delta string, opts ...EventOption) error {
	ev := NewAIMessageDelta(delta, opts...)
	return e.Emit(string(EventMessageDelta), ev, opts...)
}

// EmitMessageCitation sends "message.citation".
func (e *Emitter) EmitMessageCitation(refs []CitationReference, opts ...EventOption) error {
	ev := NewAIMessageCitation(refs, opts...)
	return e.Emit(string(EventMessageCitation), ev, opts...)
}

// EmitMessageEnd sends "message.end".
func (e *Emitter) EmitMessageEnd(reason AIEventFinishReason, opts ...EventOption) error {
	ev := NewAIMessageEnd(reason, opts...)
	return e.Emit(string(EventMessageEnd), ev, opts...)
}

// EmitError sends "error".
func (e *Emitter) EmitError(message string, opts ...EventOption) error {
	ev := NewAIError(message, opts...)
	return e.Emit(string(EventError), ev, opts...)
}
