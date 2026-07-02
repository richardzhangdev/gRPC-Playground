package main

import (
	"log"
	"net"

	pb "gRPC-Playground/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterLLMServiceServer(grpcServer, &server{})

	log.Println("gRPC server listening on :50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
