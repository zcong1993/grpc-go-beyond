package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	"github.com/zcong1993/grpc-go-beyond/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"

	_ "github.com/zcong1993/grpc-go-beyond/internal/codec"
	"github.com/zcong1993/grpc-go-beyond/internal/server"
)

func main() {
	port := flag.String("port", ":8888", "listen port")
	serverType := flag.String("type", "", "server type: stream | default")
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}

	var s *grpc.Server

	if *serverType == "stream" {
		s = grpc.NewServer(grpc.UnknownServiceHandler(server.Handler()))
	} else {
		s = grpc.NewServer()
		pb.RegisterHelloServer(s, &server.HelloServer{})
		reflection.Register(s)
	}

	service.RegisterChannelzServiceToServer(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
