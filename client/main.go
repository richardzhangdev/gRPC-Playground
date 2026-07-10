package main

import (
	"context"
	"fmt"
	"log"

	pb "gRPC-Playground/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLLMServiceClient(conn)

	req := &pb.ChatRequest{
		Model:  "deepseek",
		Prompt: "explain protobufs",
	}

	res, err := client.Chat(context.Background(), req)
	if err != nil {
		log.Fatalf("could not call Chat: %v", err)
	}

	fmt.Println("Content:", res.Content)
}
