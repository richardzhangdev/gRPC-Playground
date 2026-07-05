package types

import (
	"time"
)

type UsageEvent struct {
	RequestID        string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Latency          time.Duration
}
