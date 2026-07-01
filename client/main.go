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
	client := pb.NewGreeterClient(conn)

	// Create the request.
	req := &pb.HelloRequest{
		Name: "Richard",
	}

	// Call the remote SayHello function.
	res, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatalf("could not call SayHello: %v", err)
	}

	// Print the response.
	fmt.Println("Server replied:", res.Message)
}
