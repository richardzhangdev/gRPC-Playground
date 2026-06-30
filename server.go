package main

import (
	"context"

	"gRPC-Playground/proto"
	pb "gRPC-Playground/proto"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "Hello " + req.Name,
	}, nil
}
