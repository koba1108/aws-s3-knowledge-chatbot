package sse

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	RoleTool      Role = "tool"
)

type AIEventType string

const (
	EventMessageStart AIEventType = "message.start" // メッセージ開始（ヘッダ）
	EventMessageDelta AIEventType = "message.delta" // トークン／チャンク差分
	EventMessageEnd   AIEventType = "message.end"   // メッセージ終了（メタ／finish reason）
	EventError        AIEventType = "error"         // エラー（非ストリーム時も共通）
)

type AIEventFinishReason string

const (
	FinishCompleted     AIEventFinishReason = "completed"
	FinishStop          AIEventFinishReason = "stop"
	FinishLength        AIEventFinishReason = "length"
	FinishContentFilter AIEventFinishReason = "content_filter"
	FinishTool          AIEventFinishReason = "tool"
	FinishGuardrail     AIEventFinishReason = "guardrail_intervention"
	FinishError         AIEventFinishReason = "error"
	FinishUnknown       AIEventFinishReason = "unknown"
)

// AIBaseEvent holds fields common to all SSE events.
type AIBaseEvent struct {
	ID        string `json:"id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

type AIMessageHeader struct {
	Role Role `json:"role"`
}

type AIMessageStart struct {
	AIBaseEvent
	Type    AIEventType     `json:"type"` // "message.start"
	Message AIMessageHeader `json:"message"`
}

type AIMessageDelta struct {
	AIBaseEvent
	Type  AIEventType `json:"type"`  // "message.delta"
	Delta string      `json:"delta"` // テキスト差分
}

type AIMessageEnd struct {
	AIBaseEvent
	Type         AIEventType         `json:"type"`          // "message.end"
	FinishReason AIEventFinishReason `json:"finish_reason"` // "completed" など
}

type AIError struct {
	AIBaseEvent
	Type      AIEventType `json:"type"`
	Message   string      `json:"message"`
	Code      string      `json:"code,omitempty"`      // e.g. "AccessDenied"
	Retryable bool        `json:"retryable,omitempty"` // 再試行ヒント
}

// EventOption is a functional option to set optional fields on AIBaseEvent.
type EventOption func(*AIBaseEvent)

// WithID sets the event ID.
func WithID(id string) EventOption {
	return func(b *AIBaseEvent) { b.ID = id }
}

// WithSessionID sets the session ID.
func WithSessionID(sessionID string) EventOption {
	return func(b *AIBaseEvent) { b.SessionID = sessionID }
}

// NewAIMessageStart creates a message.start event.
// Required: role
// Optional: use EventOption (WithID, WithSessionID)
func NewAIMessageStart(role Role, opts ...EventOption) AIMessageStart {
	ev := AIMessageStart{
		AIBaseEvent: AIBaseEvent{},
		Type:        EventMessageStart,
	}
	ev.Message.Role = role
	for _, opt := range opts {
		opt(&ev.AIBaseEvent)
	}
	return ev
}

// NewAIMessageDelta creates a message.delta event.
// Required: delta
// Optional: use EventOption (WithID, WithSessionID)
func NewAIMessageDelta(delta string, opts ...EventOption) AIMessageDelta {
	ev := AIMessageDelta{
		AIBaseEvent: AIBaseEvent{},
		Type:        EventMessageDelta,
		Delta:       delta,
	}
	for _, opt := range opts {
		opt(&ev.AIBaseEvent)
	}
	return ev
}

// NewAIMessageEnd creates a message.end event.
// Required: reason
// Optional: use EventOption (WithID, WithSessionID)
func NewAIMessageEnd(reason AIEventFinishReason, opts ...EventOption) AIMessageEnd {
	ev := AIMessageEnd{
		AIBaseEvent:  AIBaseEvent{},
		Type:         EventMessageEnd,
		FinishReason: reason,
	}
	for _, opt := range opts {
		opt(&ev.AIBaseEvent)
	}
	return ev
}

// NewAIError creates an error event.
// Required: message
// Optional: use EventOption (WithID, WithSessionID)
func NewAIError(message string, opts ...EventOption) AIError {
	ev := AIError{
		AIBaseEvent: AIBaseEvent{},
		Type:        EventError,
		Message:     message,
	}
	for _, opt := range opts {
		opt(&ev.AIBaseEvent)
	}
	return ev
}

func NewAssistantStart(opts ...EventOption) AIMessageStart {
	return NewAIMessageStart(RoleAssistant, opts...)
}

func NewAssistantDelta(delta string, opts ...EventOption) AIMessageDelta {
	return NewAIMessageDelta(delta, opts...)
}

func NewAssistantEnd(reason AIEventFinishReason, opts ...EventOption) AIMessageEnd {
	return NewAIMessageEnd(reason, opts...)
}
