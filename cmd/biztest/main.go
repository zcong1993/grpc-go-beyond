package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/zcong1993/grpc-go-beyond/pb/test"
	"google.golang.org/grpc/reflection"

	_ "github.com/zcong1993/grpc-go-beyond/internal/codec"
	"google.golang.org/grpc"
)

type module struct {
	test.UnimplementedTestServer
}

func (m *module) Test(ctx context.Context, request *test.TestRequest) (*test.TestRequest, error) {
	fmt.Println("Test receive: ", request.Name)
	return request, nil
}

func (m *module) Test1(ctx context.Context, request *test.TestRequest) (*test.TestRequest, error) {
	fmt.Println("Test1 receive: ", request.Name)
	panic("implement me")
}

func main() {
	lis, err := net.Listen("tcp", ":8889")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	test.RegisterTestServer(s, &module{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
