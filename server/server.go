package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

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

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}
type Choice struct {
	Message Message `json:"message"`
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

	req_body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:4000/v1/chat/completions", "application/json", bytes.NewBuffer(req_body))
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

	var response ChatCompletionResponse
	json.Unmarshal(resp_body, &response)
	return &pb.ChatResponse{
		Content: response.Choices[0].Message.Content,
	}, nil

}
