package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/zcong1993/grpc-go-beyond/internal/proxy"
	_ "github.com/zcong1993/grpc-go-beyond/internal/proxy/codec"
)

func main() {
	port := flag.String("port", ":9999", "listen port")
	upstream := flag.String("upstream", "127.0.0.1:8888", "upstream service addr")
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnknownServiceHandler(proxy.TransparentHandler(proxy.NewManager(*upstream).StreamDirector)))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
