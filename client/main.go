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
	// Connect to the gRPC server.
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a Greeter client.
	client := pb.NewLLMServiceClient(conn)

	// Create the request.
	req := &pb.ChatRequest{
		Model:  "glm",
		Prompt: "explain protobufs",
	}

	// Call the remote SayHello function.
	res, err := client.Chat(context.Background(), req)
	if err != nil {
		log.Fatalf("could not call Chat: %v", err)
	}

	// Print the response.
	fmt.Println("Server replied:", res.Content)
}
