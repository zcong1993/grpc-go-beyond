package main

import (
	"flag"
	"log"
	"net"

	"github.com/zcong1993/grpc-go-beyond/internal/server"
	"github.com/zcong1993/grpc-go-beyond/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := flag.String("port", ":8888", "listen port")
	serverType := flag.String("type", "stream", "server type: stream | default")
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}

	var s *grpc.Server

	if *serverType == "default" {
		s = grpc.NewServer()
		pb.RegisterHelloServer(s, &server.HelloServer{})
		reflection.Register(s)
	} else {
		s = grpc.NewServer(grpc.UnknownServiceHandler(server.Handler()))
	}

	service.RegisterChannelzServiceToServer(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
