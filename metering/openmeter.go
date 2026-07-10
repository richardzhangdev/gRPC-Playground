package metering

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gRPC-Playground/types"
)

type OpenMeterClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

type OpenMeterEvent struct {
	SpecVersion string                 `json:"specversion"`
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	Type        string                 `json:"type"`
	Subject     string                 `json:"subject"`
	Time        time.Time              `json:"time"`
	Data        map[string]interface{} `json:"data"`
}

type OpenMeterResponse struct {
	EventID string `json:"id"`
}

func NewOpenMeterClient(baseURL, apiKey string) *OpenMeterClient {
	return &OpenMeterClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *OpenMeterClient) IngestEvent(ctx context.Context, usage *types.UsageEvent) error {
	event := OpenMeterEvent{
		SpecVersion: "1.0",
		ID:          usage.RequestID,
		Source:      "llm-gateway",
		Type:        "usage_event",
		Subject:     usage.Model,
		Time:        time.Now(),
		Data: map[string]interface{}{
			"model":             usage.Model,
			"prompt_tokens":     usage.PromptTokens,
			"completion_tokens": usage.CompletionTokens,
			"total_tokens":      usage.TotalTokens,
			"latency_ms":        usage.Latency.Milliseconds(),
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal failed")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/v1/events", c.baseURL),
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending event: %w", err)
	}
	defer resp.Body.Close()

	// OpenMeter returns 204 No Content on a successful ingest; accept any 2xx.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil

}
