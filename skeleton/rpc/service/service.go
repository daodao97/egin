package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "skeleton/rpc/pb"
)

// server 用来实现 hello.HelloServer

type server struct{}

// 实现 hello.SayHello 方法

// (context.Context, *HelloRequest) (*HelloReply, error)

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name, Code: 200}, nil
}

type serverB struct{}

func (s *serverB) SayHello(ctx context.Context, in *pb.ByeRequest) (*pb.ByeReply, error) {
	return &pb.ByeReply{Message: "bye " + in.Name, Code: 200}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &server{})
	pb.RegisterByeServer(s, &serverB{})

	// 在 server 中 注册 gRPC 的 reflection service

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
