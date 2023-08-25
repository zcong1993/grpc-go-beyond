package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	"github.com/zcong1993/grpc-go-beyond/pb"

	_ "github.com/zcong1993/grpc-go-beyond/internal/codec"
	"google.golang.org/grpc"
)

type module struct {
	pb.UnimplementedHelloServer
}

func (m *module) Echo(ctx context.Context, request *pb.EchoRequest) (*pb.EchoRequest, error) {
	fmt.Println("Echo receive: ", request.Message)
	return request, nil
}

func (m *module) Echo2(ctx context.Context, request *pb.EchoRequest) (*pb.EchoRequest, error) {
	fmt.Println("Echo2 receive: ", request.Message)
	return request, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &module{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
