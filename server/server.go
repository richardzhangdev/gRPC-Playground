package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	pb "gRPC-Playground/proto"
	"net/http"
)

type server struct {
	pb.UnimplementedLLMServiceServer
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type LiteLLMResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type UsageEvent struct {
	RequestID        string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Latency          time.Duration
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}

func (s *server) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	request := ChatCompletionRequest{
		Model: req.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
	}

	client := &http.Client{}

	req_body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	start_time := time.Now()
	resp, err := client.Post("http://localhost:4000/v1/chat/completions", "application/json", bytes.NewBuffer(req_body))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned status %d", resp.StatusCode)
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	latency := time.Since(start_time)

	var response LiteLLMResponse
	err = json.Unmarshal(resp_body, &response)
	if err != nil {
		return nil, err
	}

	usageEvent := UsageEvent{
		RequestID:        response.ID,
		Model:            response.Model,
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
		Latency:          latency,
	}

	s.RecordUsage(usageEvent)

	return &pb.ChatResponse{
		Content: response.Choices[0].Message.Content,
	}, nil

}

func (s *server) RecordUsage(usage UsageEvent) {
	log.Printf(
		"request_id=%s model=%s prompt_tokens=%d completion_tokens=%d total_tokens=%d latency=%s",
		usage.RequestID,
		usage.Model,
		usage.PromptTokens,
		usage.CompletionTokens,
		usage.TotalTokens,
		usage.Latency,
	)
}
