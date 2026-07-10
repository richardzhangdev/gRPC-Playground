package main

import (
	"log"
	"net"

	pb "gRPC-Playground/proto"

	"google.golang.org/grpc"

	"gRPC-Playground/consumer"
	"gRPC-Playground/metering"
	"gRPC-Playground/producer"
	"gRPC-Playground/server"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	brokers := []string{"localhost:9092"}
	topic := "usage-events"

	p := producer.NewProducer(brokers, topic)
	om := metering.NewOpenMeterClient("http://localhost:8888", "")
	go consumer.StartConsumer(brokers, topic, "usage-group", om)

	s := server.NewServer(p)

	pb.RegisterLLMServiceServer(grpcServer, s)

	log.Println("gRPC server listening on :50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
