package main

import (
	"context"

	pb "gRPC-Playground/proto"
)

type server struct {
	pb.UnimplementedLLMServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}

func (s *server) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	return &pb.ChatResponse{
		Content: "placeholder content",
	}, nil
}
